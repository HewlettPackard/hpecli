// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"encoding/json"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/spf13/cobra"
)

func init() {
	ovGetCmd.AddCommand(servers)
}

// login represents the oneview login command
var servers = &cobra.Command{
	Use:   "servers",
	Short: "Get servers from OneView: hpecli oneview get servers",
	RunE:  runGetServers,
}

func runGetServers(_ *cobra.Command, _ []string) error {
	host, apiKey := apiKey()
	ovc := NewOVClientFromAPIKey(host, apiKey)

	sh, err := ovc.GetServerHardwareList(nil, "", "", "", "")
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
