// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"github.com/spf13/cobra"
)

// NewAnalyticsCommand to turn on or off GA
func NewAnalyticsCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "analytics",
		Short: "Analytics for HPE CLI commands",
	}

	cmd.AddCommand(
		newOnCommand(),
		newOffCommand(),
		newStatusCommand(),
	)

	return cmd
}
