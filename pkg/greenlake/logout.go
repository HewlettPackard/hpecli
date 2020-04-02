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

	cmd.Flags().StringVar(&host, "host", "", "oneview host/ip address")

	return cmd
}

func runLogout(hostParam string) error {
	host, _, err := hostToLogout(hostParam)
	if err != nil {
		logrus.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for HPE GreenLake. " +
			"Please login to OneView using: hpe greenlake login")
	}


	sd, err := defaultSessionData()
	if err != nil {
		logrus.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for HPE GreenLake. " +
			"Please login to GreenLake using: hpe greenlake login")
	}

	logrus.Debugf("Attempting get greenlake users at: %v", sd.Host)


	
	// No method to logout yet

	logrus.Infof("Successfully logged out of HPE GreenLake: %s", host)

	// Cleanup context
	err = deleteSavedHostData(host)
	if err != nil {
		logrus.Warning("Unable to cleanup the session data")
		return err
	}

	return nil
}

func hostToLogout(hostParam string) (host, token string, err error) {
	if hostParam == "" {
		// they didn't specify a host.. so use the context to find one
		h, t, e := hostAndToken()
		if e != nil {
			logrus.Debugf("unable to retrieve apiKey because of: %v", e)
			return "", "", fmt.Errorf("unable to retrieve the last login for HPE GreenLake.  " +
				"Please login to OneView using: hpe greenlake login")
		}

		return h, t, nil
	}

	token, err = hostData(hostParam)
	if err != nil {
		return "", "", err
	}

	return hostParam, token, nil
}