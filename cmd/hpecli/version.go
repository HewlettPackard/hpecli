// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package main

import (
	"github.com/spf13/cobra"
	"github.hpe.com/blaine-southam/hpecli/pkg/logger"
	"github.hpe.com/blaine-southam/hpecli/pkg/version"
)

var verbose bool

// versionCmd version
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays version of hpecli",
	Run: func(_ *cobra.Command, _ []string) {
		if verbose || logger.Level >= logger.DebugLevel {
			logger.Always(version.GetFull())
		} else {
			logger.Always(version.Get())
		}
	},
}

func init() {
	versionCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose version")
}
