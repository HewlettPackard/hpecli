// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(cmdIloLogin)
}

// Cmd represents the ilo command
var Cmd = &cobra.Command{
	Use:   "ilo",
	Short: "Access to HPE iLO commands",
}
