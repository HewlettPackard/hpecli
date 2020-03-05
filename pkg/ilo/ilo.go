// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(cmdILOContext)
	Cmd.AddCommand(cmdILOGet)
	Cmd.AddCommand(cmdILOLogin)
	Cmd.AddCommand(iloLogoutCmd)
}

// Cmd represents the ilo command
var Cmd = &cobra.Command{
	Use:   "ilo",
	Short: "Access to HPE iLO commands",
}
