// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"github.com/spf13/cobra"
)

// Cmd represents the ilo command
var ovGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get data from HPE OneView",
}
