// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package cfm

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newLogoutCommand() *cobra.Command {
	var host string

	var cmd = &cobra.Command{
		Use:   "logout",
		Short: "Logout from CFM: hpecli cfm logout",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runLogout(host)
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "cfm host/ip address")

	return cmd
}

func runLogout(hostParam string) error {
	host, token, err := hostToLogout(hostParam)
	if err != nil {
		logrus.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for CFM.  " +
			"Please login to CFM using: hpecli cfm login")
	}

	cfmClient := newCFMClientFromAPIKey(host, token)

	logrus.Warningf("Using CFM: %s\n", host)

	// Use CFMClient to logout
	_, deletionError := cfmClient.DeleteAuthToken()
	if deletionError != nil {
		logrus.Warningf("Unable to logout from CFM at: %s", host)
		return errors.New(deletionError.Result)
	}

	logrus.Warningf("Successfully logged out of CFM: %s", host)

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
			return "", "", fmt.Errorf("unable to retrieve the last login for CFM.  " +
				"Please login to CFM using: hpecli cfm login")
		}

		return h, t, nil
	}

	token, err = hostData(hostParam)
	if err != nil {
		return "", "", err
	}

	return hostParam, token, nil
}
