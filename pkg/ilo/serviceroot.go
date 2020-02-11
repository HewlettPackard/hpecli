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
	context := iloContext()

	host, apiKey, err := context.APIKey()
	if err != nil {
		logger.Debug("unable to retrieve apiKey for host: %s because of: %#v", host, err)
		return fmt.Errorf("unable to retrieve the last login for HPE iLO." +
			"Please login to iLO using: hpecli ilo login")
	}
	logger.Debug("Attempting get ilo service root at: %v", host)

	c := NewILOClientFromAPIKey(host, apiKey)

	jsonResult, err := c.GetServiceRoot()
	if err != nil {
		return err
	}

	logger.Always("%s", jsonResult)

	return nil
}
