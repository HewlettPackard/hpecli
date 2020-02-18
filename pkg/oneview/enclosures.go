// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"encoding/json"
	"fmt"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/spf13/cobra"
)

var ovEnclosureData struct {
	name string
}

func init() {
	ovGetCmd.AddCommand(enclosuresCmd)
	enclosuresCmd.Flags().StringVar(&ovEnclosureData.name, "name", "", "name of the enclosure to retrieve")
}

// login represents the oneview login command
var enclosuresCmd = &cobra.Command{
	Use:   "enclosures",
	Short: "Get enclosures from OneView: hpecli oneview get enclosures",
	RunE:  getEnclosures,
}

func getEnclosures(_ *cobra.Command, _ []string) error {
	return getEnclosuresData()
}

func getEnclosuresData() error {
	d, err := getContext()
	if err != nil {
		logger.Debug("unable to retrieve apiKey because of: %#v", err)
		return fmt.Errorf("unable to retrieve the last login for OneView." +
			"Please login to OneView using: hpecli oneview login")
	}

	ovc := NewOVClientFromAPIKey(d.Host, d.APIKey)

	// not sure we want to show the host we are retieving from.
	// it's good to know - but then breaks json data format being returned
	logger.Always("Retrieving data from: %s", d.Host)

	var el interface{}

	if ovEnclosureData.name != "" {
		el, err = ovc.GetEnclosureByName(ovEnclosureData.name)
	} else {
		el, err = ovc.GetEnclosures("", "", "", "", "")
	}

	if err != nil {
		logger.Warning("Unable to login with supplied credentials to OneView at: %s", d.Host)
		return err
	}

	out, err := json.MarshalIndent(el, "", "  ")
	if err != nil {
		logger.Warning("Unable to output data as JSON.  Please try the command again.")
	}

	logger.Always("%s", out)

	return nil
}
