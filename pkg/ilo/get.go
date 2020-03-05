// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"github.com/spf13/cobra"
)

func init() {
	cmdILOGet.AddCommand(cmdILOServiceRoot)
}

// Cmd represents the ilo command
var cmdILOGet = &cobra.Command{
	Use:   "get",
	Short: "Get details from iLO",
}
