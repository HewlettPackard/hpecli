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

// GAClientIDKey - Analytics Client ID Key
const GAClientIDKey = "GA_CLIENT_ID"

// analyticsStateKey - maintain Analytics enable/disable status
// true is enabled.  false is disabled
const analyticsStateKey = "ANALYTICS_STATE"

var onOff = map[bool]string{true: "enabled", false: "disabled"}

// client - wrapper class for Google Analytics Measurement Protocol api's
type client struct {
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

// newAnalyticsClient create
func newAnalyticsClient(version, eventHitType, eventCategory, eventAction,
	eventValue, eventLabel, userAgent, applicationVersion,
	applicationName string) *client {
	return &client{
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

func SendEvent(module, command, subcommand string) {
	if !analyticsEnabled() {
		logrus.Debugf("Analytics disabled .. skipping sending event.")
	}

	client := newAnalyticsClient("1", "event", module, command,
		"200", subcommand, "hpe/0.0.1", "0.0.1", "hpecli")

	err := client.trackEvent()
	if err != nil {
		logrus.Debugf("Failure to send analytics event. Error: %+v", err)
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

// enableAnalytics enable google analytics for HPE CLI
func enableAnalytics() error {
	d, err := db.Open()
	if err != nil {
		logrus.Debug("Unable to open DB to get analyticsStateKey")

		return err
	}

	defer d.Close()

	err = d.Put(analyticsStateKey, true)
	if err != nil {
		logrus.Debug("Unable to enable analytics in DB")
	}

	return err
}

// disableAnalytics disable google analytics for HPE CLI
func disableAnalytics() error {
	d, err := db.Open()
	if err != nil {
		logrus.Debug("Unable to open DB to get analyticsStateKey")
		return err
	}

	defer d.Close()

	err = d.Put(analyticsStateKey, false)
	if err != nil {
		logrus.Debug("Unable to enable analytics in DB")
	}

	return err
}

// analyticsEnabled Check whether analytics is enabled or disabled in db
func analyticsEnabled() bool {
	d, err := db.Open()
	if err != nil {
		logrus.Debug("Unable to open DB to get analyticsStateKey")
		return false
	}

	defer d.Close()

	var enabled bool

	err = d.Get(analyticsStateKey, &enabled)
	if err != nil {
		enabled = false

		if errors.Is(err, db.ErrNotFound) {
			logrus.Debug("analyticsStateKey is not set in the DB. Defaulting to disabled")
		} else {
			logrus.Debug("Unable to determine analytics state in DB.  Defaulting to disabled.")
		}
	}

	logrus.Debugf("Analytics state: %s", onOff[enabled])

	return enabled
}

// trackEvent Measurement Protocol api to track user events
func (c *client) trackEvent() error {
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
