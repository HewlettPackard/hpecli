// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/spf13/cobra"
)

// cmdIloLogin represents the get command
var cmdILOServiceRoot = &cobra.Command{
	Use:           "serviceroot",
	Short:         "Get service root details",
	SilenceErrors: true,
	RunE:          runILOServiceRoot,
}

func runILOServiceRoot(cmd *cobra.Command, _ []string) error {
	log.Logger.Debug("Beginning runILOServiceRoot")

	sd, err := defaultSessionData()
	if err != nil {
		log.Logger.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for HPE iLO.  " +
			"Please login to iLO using: hpecli ilo login")
	}

	log.Logger.Warningf("Using iLO: %s", sd.Host)

	client := NewILOClientFromAPIKey(sd.Host, sd.Token)

	jsonResult, err := client.getServiceRoot()
	if err != nil {
		return err
	}

	log.Logger.Infof("%s", jsonResult)

	return nil
}
