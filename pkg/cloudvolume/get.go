// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"github.com/spf13/cobra"
)

func newGetCommand() *cobra.Command {
	// cmd represents the cloudvolume command
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Get resources from HPE Nimble Cloud Volumes",
	}

	cmd.AddCommand(
		newGetVolumesCommand(),
	)

	return cmd
}
