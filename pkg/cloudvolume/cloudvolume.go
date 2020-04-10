// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"github.com/spf13/cobra"
)

func NewCloudVolumeCommand() *cobra.Command {
	// cmd represents the ilo command
	var cmd = &cobra.Command{
		Use:   "cloudvolumes",
		Short: "Access to HPE Cloud Volumes commands",
	}

	cmd.AddCommand(
		newGetCommand(),
		newLoginCommand(),
		newLogoutCommand(),
	)

	return cmd
}
