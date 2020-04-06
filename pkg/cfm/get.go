// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package cfm

import (
	"github.com/spf13/cobra"
)

func newGetCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Get data from HPE Composable Fabric Manager",
	}

	cmd.AddCommand()

	return cmd
}
