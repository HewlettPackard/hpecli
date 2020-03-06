// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"encoding/json"
	"fmt"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/spf13/cobra"
)

func newEnclosuresCommand() *cobra.Command {
	var enclosureName string

	var cmd = &cobra.Command{
		Use:   "enclosures",
		Short: "Get enclosures from OneView: hpecli oneview get enclosures",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getEnclosuresData(enclosureName)
		},
	}

	cmd.Flags().StringVar(&enclosureName, "name", "", "name of the enclosure to retrieve")

	return cmd
}

func getEnclosuresData(enclosureName string) error {
	host, token, err := hostAndToken()
	if err != nil {
		log.Logger.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for OneView.  " +
			"Please login to OneView using: hpecli oneview login")
	}

	ovc := newOVClientFromAPIKey(host, token)

	log.Logger.Warningf("Using OneView: %s", host)

	var el interface{}

	if enclosureName != "" {
		el, err = ovc.GetEnclosureByName(enclosureName)
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
