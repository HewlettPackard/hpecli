// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package version

import (
	"fmt"

	"github.com/HewlettPackard/hpecli/pkg/update"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Info struct {
	Sematic     string
	GitCommitID string
	BuildDate   string
	Verbose     bool
}

func (vi *Info) String() string {
	if vi.Verbose {
		return fmt.Sprintf("%s:%s:%s", vi.Sematic, vi.GitCommitID, vi.BuildDate)
	}

	return vi.Sematic
}

func NewVersionCommand(vInfo *Info) *cobra.Command {
	var verbose bool

	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Displays version of hpecli",
		Run: func(_ *cobra.Command, _ []string) {
			vInfo.Verbose = vInfo.Verbose || verbose
			runVersion(vInfo)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose version")

	return cmd
}

func runVersion(vInfo *Info) {
	logrus.Debugln("update : starting check to see if an update is available")
	logrus.Infof("Local CLI version: %s", vInfo.String())

	resp, err := update.CheckForUpdate(vInfo.Sematic)
	if err != nil {
		logrus.Debugf("Unable to get remove version.  %s", err)
		logrus.Infoln("Unable to determine remote version.")

		return
	}

	if !resp.UpdateAvailable {
		logrus.Infoln("You are currently running the latest version of the CLI.")
		return
	}

	logrus.Infof("Updated CLI is available.  Remote version: %s.  You can update by running \"hpecli update\"",
		resp.RemoteVersion)
}
