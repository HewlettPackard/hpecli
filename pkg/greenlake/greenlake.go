// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"github.com/spf13/cobra"
)

func NewGreenlakeCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "greenlake",
<<<<<<< HEAD
		Short: "Access to HPE GreenLake commands",
=======
		Short: "Access to HPE Greenlake commands",
>>>>>>> 93c21ccd1f04dcbdbae5c791a7a7017643d16921
	}

	cmd.AddCommand(
		newGetCommand(),
		newLoginCommand(),
		newLogoutCommand(),
	)

	return cmd
}
