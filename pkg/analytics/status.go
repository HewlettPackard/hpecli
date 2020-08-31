// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newStatusCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "status",
		Short: "Show analytics state.  Enabled or Disabled",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := runAnalyticsStatus()
			SendEvent("analytics", "analytics", cmd.Name(), err)
			return err
		},
	}

	return cmd
}

func runAnalyticsStatus() error {
	enabled := analyticsEnabled()

	logrus.Infof("Anonymouns analytics reporting is %s", stateMap[enabled])

	return nil
}
