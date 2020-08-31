// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newOffCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "off",
		Short: "Turn off Analytics",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := runDisableAnalytics()
			SendEvent("analytics", "analytics", cmd.Name(), err)
			return err
		},
	}

	return cmd
}

func runDisableAnalytics() error {
	err := disableAnalytics()
	if err != nil {
		logrus.Warningf("Unable to disable analytics. Error: %+v", err)
		return err
	}

	logrus.Info("Anonymous analytics reporting has been disabled")

	return nil
}
