// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"github.com/spf13/cobra"
)

func newGetCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Get details from iLO",
	}

	cmd.AddCommand(
		newServiceRootCommand(),
	)

	return cmd
}
