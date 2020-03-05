// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"github.com/spf13/cobra"
)

// Cmd represents the ilo command
var Cmd = &cobra.Command{
	Use:   "cloudvolume",
	Short: "Access to HPE Nimble Cloud Volume commands",
}

func init() {
	Cmd.AddCommand(cmdGet)
	Cmd.AddCommand(cmdLogin)
}
