// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/spf13/cobra"
)

// cmdIloLogin represents the get command
var cmdILOServiceRoot = &cobra.Command{
	Use:   "serviceroot",
	Short: "Get service root details",
	RunE:  runILOServiceRoot,
}

func runILOServiceRoot(_ *cobra.Command, _ []string) error {
	logger.Debug("Beginning runILOServiceRoot")

	sd, err := defaultSessionData()
	if err != nil {
		logger.Debug("unable to retrieve apiKey because of: %#v", err)
		return fmt.Errorf("unable to retrieve the last login for HPE iLO." +
			"Please login to iLO using: hpecli ilo login")
	}

	logger.Debug("Attempting get ilo service root at: %v", sd.Host)

	client := NewILOClientFromAPIKey(sd.Host, sd.Token)

	jsonResult, err := client.getServiceRoot()
	if err != nil {
		return err
	}

	logger.Always("%s", jsonResult)

	return nil
}
