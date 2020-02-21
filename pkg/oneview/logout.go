// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"fmt"

	"github.com/HewlettPackard/hpecli/pkg/logger"
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
	d, err := getContext()
	if err != nil {
		logger.Debug("unable to retrieve apiKey because of: %#v", err)
		return fmt.Errorf("unable to retrieve the last login for OneView." +
			"Please login to OneView using: hpecli oneview login")
	}

	ovc := NewOVClientFromAPIKey(d.Host, d.APIKey)

	logger.Always("Retrieving data from: %s", d.Host)

	// Use OVClient to logout
	err = ovc.SessionLogout()
	if err != nil {
		logger.Warning("Unable to logout from OneView at: %s", d.Host)
		return err
	}
	
	// Cleanup context
	err = removeContext(d.Host)
	if err != nil {
		logger.Warning("Unable to cleanup the session data")
		return err
	} 

	return nil
}
