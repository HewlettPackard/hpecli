// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"github.com/spf13/cobra"
)

func newGetCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Get data from HPE OneView",
	}

	cmd.AddCommand(
		newEnclosuresCommand(),
		newServersCommand(),
	)

	return cmd
}
