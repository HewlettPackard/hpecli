// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package version

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// values are injected by the linker at build time via ldflags
var buildDate = "0"
var gitCommitID = "0"
var version = "0.0.0"

func NewVersionCommand() *cobra.Command {
	var verbose bool

	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Displays version of hpecli",
		Run: func(_ *cobra.Command, _ []string) {
			runVersion(verbose)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose version")

	return cmd
}

func runVersion(verbose bool) {
	logrus.Info(versionToShow(verbose))
}

func versionToShow(verbose bool) string {
	ver := Get()
	if isFullVersion(verbose) {
		ver = GetFull()
	}

	return ver
}

func isFullVersion(verbose bool) bool {
	return verbose || logrus.GetLevel() >= logrus.DebugLevel
}

// Get returns the short version. just the version (e.g. 0.0.1)
func Get() string {
	return fmt.Sprintf("%s\n", version)
}

// GetFull returns the long version. (e.g. 0.0.2:6683f37:2019-11-23)
func GetFull() string {
	return fmt.Sprintf("%s:%s:%s\n", version, gitCommitID, buildDate)
}
