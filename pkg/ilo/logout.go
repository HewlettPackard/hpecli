// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/spf13/cobra"
)

var iloLogoutHost struct {
	host string
}

// iloLogoutCmd represents the ilo logout command
var iloLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from ilo: hpecli ilo logout",
	RunE:  runILOLogout,
}

func init() {
	iloLogoutCmd.Flags().StringVar(&iloLogoutHost.host, "host", "", "ilo host/ip address")
}

func runILOLogout(_ *cobra.Command, _ []string) error {
	log.Logger.Debug("Beginning runILOLogout")

	sessionData, err := sessionDataToLogout()
	if err != nil {
		log.Logger.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for HPE iLO.  " +
			"Please login to iLO using: hpecli ilo login")
	}

	log.Logger.Debugf("Attempting get ilo service root at: %v", sessionData.Host)

	client := NewILOClientFromAPIKey(sessionData.Host, sessionData.Token)

	err = client.logout(sessionData.Location)
	if err != nil {
		log.Logger.Warningf("Unable to logout from iLO at: %s", sessionData.Host)
		return err
	}

	// Cleanup context
	err = deleteSessionData(sessionData.Host)
	if err != nil {
		log.Logger.Warning("Unable to cleanup the session data")
		return err
	}

	return nil
}

func sessionDataToLogout() (data *sessionData, err error) {
	data = &sessionData{}

	if iloLogoutHost.host == "" {
		// they didn't specify a host.. so use the context to find one
		d, e := defaultSessionData()
		if e != nil {
			log.Logger.Debugf("unable to retrieve apiKey because of: %v", e)
			return data, fmt.Errorf("unable to retrieve the last login for iLO.  " +
				"Please login to iLO using: hpecli ilo login")
		}

		return d, nil
	}

	// they specified a host to logout.  get the token for that host
	if !strings.HasPrefix(iloLogoutHost.host, "http") {
		iloLogoutHost.host = fmt.Sprintf("https://%s", iloLogoutHost.host)
	}

	d, err := getSessionData(iloLogoutHost.host)
	if err != nil {
		return data, err
	}

	return d, nil
}
