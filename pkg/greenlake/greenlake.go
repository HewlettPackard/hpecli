// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(cmdGreenLakeLogin)
	Cmd.AddCommand(cmdGreenLakeGet)
}

// Cmd represents the greenlake command
var Cmd = &cobra.Command{
	Use:   "greenlake",
	Short: "Access to HPE greenlake commands",
}
