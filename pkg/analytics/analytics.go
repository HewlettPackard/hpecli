// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type analyticsOptions struct {
	on  string
	off string
}

// NewAnalyticsCommand to turn on or off GA
func NewAnalyticsCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "analytics",
		Short: "Google Analytics for HPE CLI commands",
	}

	cmd.AddCommand(
		newOnCommand(),
		newOffCommand(),
		newStatusCommand(),
	)

	return cmd
}

func init() {
	// Create a client ID to track user for Google Analytics
	_, err := clientID()
	if err != nil {
		logrus.Debug("Error generating client ID for Google Analytics")
	}
}
