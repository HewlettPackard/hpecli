// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"github.com/spf13/cobra"
)

func NewGreenlakeCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "greenlake",
		Short: "Access to HPE Greenlake commands",
	}

	cmd.AddCommand(
		newContextCommand(),
		newGetCommand(),
		newLoginCommand(),
	)

	return cmd
}
