// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"github.com/spf13/cobra"
)

func NewCloudVolumeCommand() *cobra.Command {
	// cmd represents the ilo command
	var cmd = &cobra.Command{
		Use:   "cloudvolume",
		Short: "Access to HPE Nimble Cloud Volume commands",
	}

	cmd.AddCommand(
		newGetCommand(),
		newLoginCommand(),
	)

	return cmd
}
