// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newLogoutCommand() *cobra.Command {
	var host string

	var cmd = &cobra.Command{
		Use:   "logout",
		Short: "Logout from HPE GreenLake",
		RunE: func(_ *cobra.Command, _ []string) error {
			if host != "" && !strings.HasPrefix(host, "http") {
				host = fmt.Sprintf("https://%s", host)
			}

			return runLogout(host)
		},
	}

	cmd.Flags().StringVar(&host, "host", greenlakeDefaultHost, "oneview host/ip address")

	return cmd
}

func runLogout(host string) error {
	sessionData, err := sessionDataToLogout(host)
	if err != nil {
		logrus.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for HPE GreenLake. " +
			"Please login to HPE GreenLake using: hpe greenlake login")
	}

	// No method to logout yet
	logrus.Infof("Successfully logged out of HPE GreenLake: %s", sessionData.Host)

	// Cleanup context
	err = deleteSavedHostData(sessionData.Host)
	if err != nil {
		logrus.Warning("Unable to cleanup the session data")
		return err
	}

	return nil
}

func sessionDataToLogout(host string) (sd *sessionData, err error) {
	sd = &sessionData{}

	if host == "" {
		// they didn't specify a host.. so use the context to find one
		d, e := defaultSessionData()
		if e != nil {
			logrus.Debugf("unable to retrieve apiKey because of: %v", e)
			return sd, fmt.Errorf("unable to retrieve the last login for HPE GreenLake.  " +
				"Please login to HPE GreenLake using: hpe greenlake login")
		}

		return d, nil
	}

	d, err := getSessionData(host)
	if err != nil {
		return sd, err
	}

	return d, nil
}