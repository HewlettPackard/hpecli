// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"github.com/spf13/cobra"
)

func NewILOCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "ilo",
		Short: "Access to HPE iLO commands",
	}

	cmd.AddCommand(
		newContextCommand(),
		newGetCommand(),
		newLoginCommand(),
		newLogoutCommand(),
	)

	return cmd
}
