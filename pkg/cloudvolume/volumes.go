// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newGetVolumesCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "volumes",
		Short: "Get from Cloud Volumes: hpecli cloudvolumes get volumes",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runGetVolumes()
		},
	}

	return cmd
}

func runGetVolumes() error {
	logrus.Debug("Beginning runCVGetVolumes")

	host, token, err := hostAndToken()
	if err != nil {
		logrus.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for HPE CloudVolumes.  " +
			"Please login to CloudVolumes using: hpecli cloudvolume login")
	}

	logrus.Debugf("Attempting get cloud volumes at: %v", host)

	cvc := newCVClientFromAPIKey(host, token)

	jsonResult, err := cvc.getCloudVolumes()
	if err != nil {
		return err
	}

	logrus.Infof("%s", jsonResult)

	return nil
}
