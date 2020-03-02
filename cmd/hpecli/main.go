// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package main

import (
	"fmt"
	"os"

	"github.com/HewlettPackard/hpecli/pkg/bashcompletion"
	"github.com/HewlettPackard/hpecli/pkg/cloudvolume"
	"github.com/HewlettPackard/hpecli/pkg/greenlake"
	"github.com/HewlettPackard/hpecli/pkg/ilo"
	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/HewlettPackard/hpecli/pkg/oneview"
	"github.com/HewlettPackard/hpecli/pkg/update"
	"github.com/HewlettPackard/hpecli/pkg/version"
	"github.com/spf13/cobra"
)

func main() {
	updateAvailableChan := make(chan bool)

	go func() {
		updateAvailableChan <- update.IsUpdateAvailable()
	}()

	var rootCmd = &cobra.Command{
		Use:   "hpecli",
		Short: "hpe cli for accessing various services",
	}

	addCommands(rootCmd)

	rootCmd.SetVersionTemplate("something out")

	logLevel := rootCmd.PersistentFlags().StringP("loglevel", "l", "warning",
		"set log level.  Possible values are: debug, info, warning, critical")

	cobra.OnInitialize(func() {
		logger.Color = true
		logger.SetLogLevel(*logLevel)
	})

	const exitError = 1

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(exitError)
	}

	newRelease := <-updateAvailableChan
	if newRelease {
		logger.Always("  An updated version of the CLI is available")
	}

	//rootCmd.GenBashCompletion(os.Stdout)
}

func addCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(cloudvolume.Cmd)
	rootCmd.AddCommand(ilo.Cmd)
	rootCmd.AddCommand(oneview.Cmd)
	rootCmd.AddCommand(update.Cmd)
	rootCmd.AddCommand(greenlake.Cmd)
	rootCmd.AddCommand(version.Cmd)
	rootCmd.AddCommand(bashcompletion.Cmd)
}
