// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"fmt"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/spf13/cobra"
)

// Cmd represents the ilo command
var cmdGetVolumes = &cobra.Command{
	Use:   "cloudvolumes",
	Short: "Get from Cloud Volumes: hpecli cloudvolumes get cloudvolumes",
	RunE:  runCVGetVolumes,
}

func runCVGetVolumes(_ *cobra.Command, _ []string) error {
	log.Logger.Debug("Beginning runCVGetVolumes")

	host, token, err := hostAndToken()
	if err != nil {
		log.Logger.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for HPE CloudVolumes.  " +
			"Please login to CloudVolumes using: hpecli cloudvolume login")
	}

	log.Logger.Debugf("Attempting get cloud volumes at: %v", host)

	cvc := NewCVClientFromAPIKey(host, token)

	jsonResult, err := cvc.GetCloudVolumes()
	if err != nil {
		return err
	}

	log.Logger.Infof("%s", jsonResult)

	return nil
}
