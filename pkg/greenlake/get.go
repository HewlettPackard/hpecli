// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"github.com/spf13/cobra"
)

func newGetCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Get details from HPE GreenLake",
	}

	cmd.AddCommand(
		newUsersCommand(),
	)

	return cmd
}
