// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const analytics = "analytics"

func newOnCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "on",
		Short: "Turn on Analytics",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := runEnableAnalytics()
			SendEvent("analytics", "analytics", cmd.Name(), err)
			return err
		},
	}

	return cmd
}

func runEnableAnalytics() error {
	err := enableAnalytics()
	if err != nil {
		logrus.Warningf("Unable to enable analytics. Error: %+v", err)
		return err
	}

	logrus.Info("Anonymous analytics reporting has been enabled")

	return nil
}
