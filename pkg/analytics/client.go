// (C) Copyright 2019 Hewlett Packard Enterprise Development LP

package analytics

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/HewlettPackard/hpecli/internal/platform/db"
	"github.com/HewlettPackard/hpecli/internal/platform/rest"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const TrackingID = "UA-159515478-1"
const ClientIDKey = "GA_CLIENT_ID"

// AnalyticsClient - wrapper class for Google Analytics Measurement Protocol api's
type AnalyticsClient struct {
	Version            string
	EventHitType       string
	Eventcategory      string
	EventAction        string
	EventValue         string
	EventLabel         string
	UserAgent          string
	ApplicationName    string
	ApplicationVersion string
	*rest.Request
}

// NewAnalyticsClient create
func NewAnalyticsClient(version, eventHitType, eventCategory, eventAction, eventValue, eventLabel, userAgent, applicationVersion, applicationName string) *AnalyticsClient {
	return &AnalyticsClient{
		Version:            version,
		EventHitType:       eventHitType,
		Eventcategory:      eventCategory,
		EventAction:        eventAction,
		EventLabel:         eventLabel,
		EventValue:         eventValue,
		UserAgent:          userAgent,
		ApplicationVersion: applicationVersion,
		ApplicationName:    applicationName,
	}
}

func clientID() string {
	d, err := db.Open()
	if err != nil {
		logrus.Debug("Unable to open DB to get client ID")
		return newClientID()
	}

	defer d.Close()

	var ID string
	if err := d.Get(ClientIDKey, &ID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			// couldn't find existing Id, so generate
			// one and store it
			ID = newClientID()
			logrus.Debugf("Didn't find existing clientID, generating a new one: %s", ID)
			_ = d.Put(ClientIDKey, ID)

			return ID
		}
		// unknown error getting key
		return newClientID()
	}

	return ID
}

func newClientID() string {
	return uuid.New().String()
}

// TrackEvent Measurement Protocol api
func (c *AnalyticsClient) TrackEvent() (string, error) {

	const uriPath = "/collect"
	const host = "https://www.google-analytics.com"

	// analyticsSON := fmt.Sprintf(`{"v":"1", "tid":"%s", "t":"%s,
	// "ea":"%s", "ev":"%s", "ec":"%s", "ua":"%s", "an":"%s", "av":"%s", "cid":"%s"}`,
	// 	c.TrackingID, c.EventHitType, c.EventAction, c.EventValue, c.Eventcategory, c.UserAgent, c.ApplicationName, c.ApplicationVersion, c.ClientID)

	u, err := url.Parse("https://www.google-analytics.com/collect")
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "https"
	u.Host = "google-analytics.com"
	u.Path = "/collect"

	q := u.Query()
	q.Set("v", "1")
	q.Set("tid", TrackingID)
	q.Set("t", c.EventHitType)
	q.Set("ea", c.EventAction)
	q.Set("ev", c.EventValue)
	q.Set("ec", c.Eventcategory)
	q.Set("ua", c.UserAgent)
	q.Set("an", c.ApplicationName)
	q.Set("av", c.ApplicationVersion)
	q.Set("cid", clientID())
	u.RawQuery = q.Encode()
	fmt.Println(u)
	print("url", host+uriPath)

	resp, err := rest.Post(u.String(), nil)
	print("resp", resp.StatusCode)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to create login sessions to Green Lake.  Repsponse was: %+v", resp.Status)
	}

	return "success", nil

}
