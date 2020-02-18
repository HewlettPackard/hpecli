// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"github.com/spf13/cobra"
)

// cmdGet represents the greenlake command
var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "Get resources from HPE Green Lake",
}

func init() {
	cmdGet.AddCommand(cmdGetUsers)
}
