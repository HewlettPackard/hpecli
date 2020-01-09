// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package main

import (
	"fmt"
	"os"

	"github.com/HewlettPackard/hpecli/pkg/ilo"
	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/HewlettPackard/hpecli/pkg/update"
	"github.com/HewlettPackard/hpecli/pkg/cloudvolumes"
	
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "hpecli",
		Short: "hpe cli for accessing various services",
	}
	addCommands(rootCmd)

	logLevel := rootCmd.PersistentFlags().StringP("loglevel", "l", "warning", "set log level.  Possible values are: debug, info, warning, critical")
	cobra.OnInitialize(func() {
		logger.Color = true
		logger.SetLogLevel(*logLevel)
		if update.IsUpdateAvailable() {
			logger.Always("  An updated version of the CLI is available")
		}
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(ilo.Cmd)
	rootCmd.AddCommand(cloudvolumes.Cmd)
	rootCmd.AddCommand(versionCmd)
}
