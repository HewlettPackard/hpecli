// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package version

import (
	"fmt"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// values are injected by the linker at build time via ldflags
var buildDate = "0"
var gitCommitID = "0"
var version = "0.0.0"

var verbose bool

// Cmd version
var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Displays version of hpecli",
	Run:   run,
}

func init() {
	Cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose version")
}

func run(_ *cobra.Command, _ []string) {
	v := versionOutput()
	log.Logger.Info(v)
}

func versionOutput() string {
	if isFullVersion() {
		return GetFull()
	}

	return Get()
}

func isFullVersion() bool {
	return verbose || log.Logger.Level >= logrus.DebugLevel
}

// Get returns the short version. just the version (e.g. 0.0.1)
func Get() string {
	return version
}

// GetFull returns the long version. (e.g. 0.0.2:6683f37:2019-11-23)
func GetFull() string {
	return fmt.Sprintf("%s:%s:%s", version, gitCommitID, buildDate)
}
