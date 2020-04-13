// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cfm

import (
	"github.com/spf13/cobra"
)

func NewCFMCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "cfm",
		Short: "Access to HPE Composable Fabric commands",
	}

	cmd.AddCommand(
		newContextCommand(),
		newGetCommand(),
		newLoginCommand(),
		newLogoutCommand(),
	)

	return cmd
}
