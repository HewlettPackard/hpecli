// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"encoding/json"
	"fmt"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/spf13/cobra"
)

var ovServersData struct {
	name string
}

func init() {
	ovGetCmd.AddCommand(serversCmd)
	serversCmd.Flags().StringVar(&ovServersData.name, "name", "", "name of the server to retrieve")
}

// login represents the oneview login command
var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "Get servers from OneView: hpecli oneview get servers",
	RunE:  getServers,
}

func getServers(_ *cobra.Command, _ []string) error {
	return getServerHardware()
}

func getServerHardware() error {
	host, apiKey := apiKey()
	if apiKey == "" {
		logger.Debug("apiKey for host: %s not set", host)
		return fmt.Errorf("unable to retrieve the last login for OneView." +
			"Please login to OneView using: hpecli login OneView")
	}

	ovc := NewOVClientFromAPIKey(host, apiKey)

	logger.Always("Retrieving data from: %s", host)

	var (
		sh  interface{}
		err error
	)

	if ovServersData.name != "" {
		sh, err = ovc.GetServerHardwareByName(ovServersData.name)
	} else {
		sh, err = ovc.GetServerHardwareList(nil, "", "", "", "")
	}

	if err != nil {
		logger.Warning("Unable to login with supplied credentials to OneView at: %s", host)
		return err
	}

	out, err := json.MarshalIndent(sh, "", "  ")
	if err != nil {
		logger.Warning("Unable to output data as JSON.  Please try the command again.")
	}

	logger.Always("%s", out)

	return nil
}