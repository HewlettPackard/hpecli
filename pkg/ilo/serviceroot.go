// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/spf13/cobra"
)

func newServiceRootCommand() *cobra.Command {
	// cmd represents the ilo command
	var cmd = &cobra.Command{
		Use:           "serviceroot",
		Short:         "Get service root details",
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runILOServiceRoot()
		},
	}

	return cmd
}

func runILOServiceRoot() error {
	log.Logger.Debug("Beginning runILOServiceRoot")

	sd, err := defaultSessionData()
	if err != nil {
		log.Logger.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for HPE iLO.  " +
			"Please login to iLO using: hpecli ilo login")
	}

	log.Logger.Warningf("Using iLO: %s", sd.Host)

	client := newILOClientFromAPIKey(sd.Host, sd.Token)

	jsonResult, err := client.getServiceRoot()
	if err != nil {
		return err
	}

	log.Logger.Infof("%s", jsonResult)

	return nil
}
