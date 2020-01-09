// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolumes

import (
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(cmdCloudVolumesLogin)
	Cmd.AddCommand(cmdCloudVolumesGet)
}

// Cmd represents the cloud volumes command
var Cmd = &cobra.Command{
	Use:   "cloudvolumes",
	Short: "Access to HPE Cloud Volumes commands",
}
