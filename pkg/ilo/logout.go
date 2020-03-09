// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

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
		Short: "Logout from ilo: hpecli ilo logout",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !strings.HasPrefix(host, "http") {
				host = fmt.Sprintf("https://%s", host)
			}

			return runLogout(host)
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "ilo host/ip address")

	return cmd
}

func runLogout(host string) error {
	logrus.Debug("Beginning runILOLogout")

	sessionData, err := sessionDataToLogout(host)
	if err != nil {
		logrus.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for HPE iLO.  " +
			"Please login to iLO using: hpecli ilo login")
	}

	logrus.Warningf("Using iLO: %s\n", sessionData.Host)

	client := newILOClientFromAPIKey(sessionData.Host, sessionData.Token)

	err = client.logout(sessionData.Location)
	if err != nil {
		logrus.Warningf("Unable to logout from iLO at: %s", sessionData.Host)
		return err
	}

	logrus.Warningf("Successfully logged out of remote ilo: %s", sessionData.Host)

	// Cleanup context
	err = deleteSessionData(sessionData.Host)
	if err != nil {
		logrus.Warning("Unable to cleanup the session data")
		return err
	}

	return nil
}

func sessionDataToLogout(host string) (data *sessionData, err error) {
	data = &sessionData{}

	if host == "" {
		// they didn't specify a host.. so use the context to find one
		d, e := defaultSessionData()
		if e != nil {
			logrus.Debugf("unable to retrieve apiKey because of: %v", e)
			return data, fmt.Errorf("unable to retrieve the last login for iLO.  " +
				"Please login to iLO using: hpecli ilo login")
		}

		return d, nil
	}

	d, err := getSessionData(host)
	if err != nil {
		return data, err
	}

	return d, nil
}
