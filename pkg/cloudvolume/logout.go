// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

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
		Short: "Logout from HPE Cloud Volumes",
		RunE: func(cmd *cobra.Command, args []string) error {
			if host != "" && !strings.HasPrefix(host, "http") {
				host = fmt.Sprintf("https://%s", host)
			}
			
			return runLogout(host)
		},
	}

	cmd.Flags().StringVar(&host, "host", cvDefaultHost, "HPE Cloud Volumes portal hostname/ip")

	return cmd
}

func runLogout(host string) error {
	logrus.Debug("Beginning runCloudVolumeLogout")
	
	host, token, err := hostAndToken()
	if err != nil {
		logrus.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("Unable to retrieve the last login for HPE Cloud volumes. " +
			"Please login to HPE Cloud Volumes using: hpe cloudvolumes login")
	}

	_ = newCVClientFromAPIKey(host, token)

	// There is no API logout we can use
	logrus.Infof("Successfully logged out of HPE Cloud Volumes")

	// Cleanup context
	err = deleteSavedHostData(host)
	if err != nil {
		logrus.Warning("Unable to cleanup the session data")
		return err
	}

	return nil
}

