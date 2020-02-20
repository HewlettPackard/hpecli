// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(cmdGLContext)
	Cmd.AddCommand(cmdGLGet)
	Cmd.AddCommand(cmdGLLogin)
}

// Cmd represents the greenlake command
var Cmd = &cobra.Command{
	Use:   "greenlake",
	Short: "Access to HPE Green Lake commands",
}
