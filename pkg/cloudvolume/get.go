// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"github.com/spf13/cobra"
)

// Cmd represents the ilo command
var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "Get resources from HPE Nimble Cloud Volumes",
}

func init() {
	cmdGet.AddCommand(cmdGetVolumes)
}
