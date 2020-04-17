// (C) Copyright 2020 Hewlett Packard Enterprise Development LP

package analytics

import (
	"errors"
	"log"
	"net/url"

	"github.com/HewlettPackard/hpecli/internal/platform/db"
	"github.com/HewlettPackard/hpecli/internal/platform/rest"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// TrackingID - Google Analytics tracking ID
const TrackingID = "UA-159515478-1"

// GAClientIDKey - Google Analytics Client ID Key
const GAClientIDKey = "GA_CLIENT_ID"

// DisableAnalyticsKey - maintain Google Analytics enable/disable status
const DisableAnalyticsKey = "GA_DISABLE"

const dbPutErr = "Unable to put the key %s in to DB"

// Client - wrapper class for Google Analytics Measurement Protocol api's
type Client struct {
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
func NewAnalyticsClient(version, eventHitType, eventCategory, eventAction,
	eventValue, eventLabel, userAgent, applicationVersion,
	applicationName string) *Client {
	return &Client{
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

func clientID() (string, error) {
	d, err := db.Open()
	if err != nil {
		logrus.Debug("Unable to open DB to get client ID")
		return "", err
	}

	defer d.Close()

	var ID string
	if err := d.Get(GAClientIDKey, &ID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			// couldn't find existing Id, so generate one and store it
			ID = newClientID()
			logrus.Debugf("Didn't find existing clientID, generating a new one: %s", ID)
			_ = d.Put(GAClientIDKey, ID)

			return ID, nil
		}
		// unknown error getting key
		return "", err
	}

	return ID, nil
}

// Generates a new unique id
func newClientID() string {
	return uuid.New().String()
}

// enableGoogleAnalytics enable google analytics for HPE CLI
func enableGoogleAnalytics() (bool, error) {
	var disableGA bool

	d, err := db.Open()

	if err != nil {
		logrus.Debug("Unable to open DB to get DisableAnalyticsKey")

		return false, err
	}
	defer d.Close()

	if err := d.Get(DisableAnalyticsKey, &disableGA); err != nil {
		return true, nil
	}

	if disableGA {
		logrus.Debugf("Found existing DisableAnalyticsKey, updating it's value to : %t", disableGA)

		err := d.Put(DisableAnalyticsKey, false)
		if err != nil {
			logrus.Debugf(dbPutErr, DisableAnalyticsKey)
			return false, err
		}

		return true, nil
	}

	return true, nil
}

// disableGoogleAnalytics disable google analytics for HPE CLI
func disableGoogleAnalytics() (bool, error) {
	var disableGA bool

	d, err := db.Open()

	if err != nil {
		logrus.Debug("Unable to open DB to get GAClientIDKey ")
		return false, err
	}
	defer d.Close()

	if err := d.Get(DisableAnalyticsKey, &disableGA); err != nil {
		logrus.Debugf("Didn't find existing DisableAnalyticsKey, creating a new one: %s", DisableAnalyticsKey)

		err := d.Put(DisableAnalyticsKey, true)
		if err != nil {
			logrus.Debugf(dbPutErr, DisableAnalyticsKey)
			return false, err
		}

		return true, nil
	}

	if !disableGA {
		logrus.Debugf("Didn't find existing enableAnalytics, creating a new one: %s", DisableAnalyticsKey)

		err := d.Put(DisableAnalyticsKey, true)
		if err != nil {
			logrus.Debugf(dbPutErr, DisableAnalyticsKey)
			return false, err
		}

		return true, nil
	}

	return true, nil
}

// CheckGoogleAnalytics Check whether google analytics is enabled or disabled in db
func CheckGoogleAnalytics() (bool, error) {
	var disableGA bool

	d, err := db.Open()

	if err != nil {
		logrus.Debug("Unable to open DB to get GAClientIDKey ")
		return false, err
	}
	defer d.Close()

	if err := d.Get(DisableAnalyticsKey, &disableGA); err != nil {
		logrus.Debugf("Didn't found existing DisableAnalyticsKey key %t", disableGA)

		return true, nil
	}

	if disableGA {
		return false, nil
	}

	return true, nil
}

// TrackEvent Measurement Protocol api to track user events
func (c *Client) TrackEvent() error {
	id, err := clientID()
	if err != nil {
		logrus.Debug("error generating client ID")
		return err
	}

	u, err := url.Parse("https://www.google-analytics.com/collect")
	if err != nil {
		log.Fatal(err)
	}

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
	q.Set("cid", id)
	u.RawQuery = q.Encode()

	_, err = rest.Post(u.String(), nil)
	if err != nil {
		return err
	}

	return nil
}
