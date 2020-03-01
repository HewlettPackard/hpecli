// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/spf13/cobra"
)

var ovLogoutHost struct {
	host string
}

// ovLogoutCmd represents the oneview ovLoginCmd command
var ovLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from OneView: hpecli oneview logout",
	RunE:  runOVLogout,
}

func init() {
	ovLogoutCmd.Flags().StringVar(&ovLogoutHost.host, "host", "", "oneview host/ip address")
}

func runOVLogout(_ *cobra.Command, _ []string) error {
	host, token, err := hostToLogout()
	if err != nil {
		log.Logger.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for OneView.  " +
			"Please login to OneView using: hpecli oneview login")
	}

	ovc := NewOVClientFromAPIKey(host, token)

	log.Logger.Infof("Retrieving data from: %s", host)

	// Use OVClient to logout
	err = ovc.SessionLogout()
	if err != nil {
		log.Logger.Warningf("Unable to logout from OneView at: %s", host)
		return err
	}

	// Cleanup context
	err = deleteSavedHostData(host)
	if err != nil {
		log.Logger.Warning("Unable to cleanup the session data")
		return err
	}

	return nil
}

func hostToLogout() (host, token string, err error) {
	if ovLogoutHost.host == "" {
		// they didn't specify a host.. so use the context to find one
		h, t, e := hostAndToken()
		if e != nil {
			log.Logger.Debugf("unable to retrieve apiKey because of: %v", e)
			return "", "", fmt.Errorf("unable to retrieve the last login for OneView.  " +
				"Please login to OneView using: hpecli oneview login")
		}

		return h, t, nil
	}

	// they specified a host to logout.  get the token for that host
	if !strings.HasPrefix(ovLogoutHost.host, "http") {
		ovLogoutHost.host = fmt.Sprintf("https://%s", ovLogoutHost.host)
	}

	token, err = hostData(ovLogoutHost.host)
	if err != nil {
		return "", "", err
	}

	return ovLogoutHost.host, token, nil
}
