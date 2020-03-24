// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"github.com/spf13/cobra"
)

type analyticsOptions struct {
	on  string
	off string
}

// NewAnalyticsCommand to turn on or off GA
func NewAnalyticsCommand() *cobra.Command {
	// var opts analyticsOptions

	var cmd = &cobra.Command{
		Use:   "analytics",
		Short: "Google Analytics for HPE CLI commands",
	}
	cmd.AddCommand(
		newOnCommand(),
		newOffCommand(),
	)
	return cmd
}
