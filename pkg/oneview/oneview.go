// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"github.com/spf13/cobra"
)

// Cmd represents the ilo command
var Cmd = &cobra.Command{
	Use:   "oneview",
	Short: "Access to HPE OneView commands",
}

func init() {
	Cmd.AddCommand(
		newContextCommand(),
		newGetCommand(),
		newLoginCommand(),
		newLogoutCommand(),
	)
}
