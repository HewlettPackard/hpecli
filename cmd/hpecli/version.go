// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package main

import (
	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/HewlettPackard/hpecli/pkg/version"
	"github.com/spf13/cobra"
)

var verbose bool

// versionCmd version
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays version of hpecli",
	Run:   run,
}

func init() {
	versionCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose version")
}

func run(_ *cobra.Command, _ []string) {
	v := versionOutput()
	logger.Always(v)
}

func versionOutput() string {
	if isFullVersion() {
		return version.GetFull()
	}
	return version.Get()
}

func isFullVersion() bool {
	return verbose || logger.Level >= logger.DebugLevel
}
