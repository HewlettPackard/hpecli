// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"encoding/json"
	"fmt"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
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
	host, token, err := hostAndToken()
	if err != nil {
		log.Logger.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for OneView.  " +
			"Please login to OneView using: hpecli oneview login")
	}

	ovc := NewOVClientFromAPIKey(host, token)

	// not sure we want to show the host we are retieving from.
	// it's good to know - but then breaks json data format being returned
	log.Logger.Warningf("Retrieving data from: %s\n", host)

	var el interface{}

	if ovEnclosureData.name != "" {
		el, err = ovc.GetEnclosureByName(ovEnclosureData.name)
	} else {
		el, err = ovc.GetEnclosures("", "", "", "", "")
	}

	if err != nil {
		log.Logger.Warningf("Unable to login with supplied credentials to OneView at: %s", host)
		return err
	}

	out, err := json.MarshalIndent(el, "", "  ")
	if err != nil {
		log.Logger.Warning("Unable to output data as JSON.  Please try the command again.")
	}

	log.Logger.Infof("%s", out)

	return nil
}
