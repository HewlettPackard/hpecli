// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"fmt"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/spf13/cobra"
)

// Cmd represents the ilo command
var cmdGetVolumes = &cobra.Command{
	Use:   "cloudvolumes",
	Short: "Get from Cloud Volumes: hpecli cloudvolumes get cloudvolumes",
	RunE:  runCVGetVolumes,
}

func runCVGetVolumes(_ *cobra.Command, _ []string) error {
	logger.Debug("Attempting get cloud volumes with user: %v, at: %v", cvLoginData.username, cvLoginData.host)

	c := cvContext()

	host, apiKey, err := c.APIKey()
	if err != nil {
		logger.Debug("unable to retrieve apiKey for host: %s because of: %#v", host, err)
		return fmt.Errorf("unable to retrieve the last login for HPE CloudVolumes." +
			"Please login to CloudVolumes using: hpecli cloudvolume get cloudvolumes")
	}

	cvc := NewCVClientFromAPIKey(host, apiKey)

	jsonResult, err := cvc.GetCloudVolumes()
	if err != nil {
		return err
	}

	logger.Always("%s", jsonResult)

	return nil
}
