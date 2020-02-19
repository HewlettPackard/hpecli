// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"github.com/spf13/cobra"
)

func init() {
	cmdGLGet.AddCommand(cmdGetUsers)
}

// cmdGLGet represents the greenlake command
var cmdGLGet = &cobra.Command{
	Use:   "get",
	Short: "Get details from HPE Green Lake",
}
