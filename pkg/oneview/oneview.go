// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"github.com/spf13/cobra"
)

func NewOneViewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "oneview",
		Short: "Access to HPE OneView commands",
	}

	cmd.AddCommand(
		newContextCommand(),
		newGetCommand(),
		newLoginCommand(),
		newLogoutCommand(),
	)

	return cmd
}
