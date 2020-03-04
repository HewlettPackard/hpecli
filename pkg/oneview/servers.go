// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"encoding/json"
	"fmt"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/spf13/cobra"
)

var ovServersData struct {
	name string
}

// login represents the oneview login command
var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "Get servers from OneView: hpecli oneview get servers",
	RunE:  getServers,
}

func init() {
	ovGetCmd.AddCommand(serversCmd)
	serversCmd.Flags().StringVar(&ovServersData.name, "name", "", "name of the server to retrieve")
}

func getServers(_ *cobra.Command, _ []string) error {
	return getServerHardware()
}

func getServerHardware() error {
	host, token, err := hostAndToken()
	if err != nil {
		log.Logger.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for OneView.  " +
			"Please login to OneView using: hpecli oneview login")
	}

	ovc := NewOVClientFromAPIKey(host, token)

	log.Logger.Warningf("Using OneView: %s", host)

	var sh interface{}
	if ovServersData.name != "" {
		sh, err = ovc.GetServerHardwareByName(ovServersData.name)
	} else {
		sh, err = ovc.GetServerHardwareList(nil, "", "", "", "")
	}

	if err != nil {
		log.Logger.Warningf("Unable to login with supplied credentials to OneView at: %s", host)
		return err
	}

	out, err := json.MarshalIndent(sh, "", "  ")
	if err != nil {
		log.Logger.Warning("Unable to output data as JSON.  Please try the command again.")
	}

	log.Logger.Infof("%s", out)

	return nil
}
