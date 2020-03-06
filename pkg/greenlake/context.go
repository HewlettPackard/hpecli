// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type glContextOptions struct {
	host     string
	tenantid string
}

func newContextCommand() *cobra.Command {
	opts := &glContextOptions{}

	var cmd = &cobra.Command{
		Use:   "context",
		Short: "Change context to different GreenLake host",
		RunE: func(_ *cobra.Command, _ []string) error {
			if !strings.HasPrefix(opts.host, "http") {
				opts.host = fmt.Sprintf("https://%s", opts.host)
			}

			return setContext(opts.host)
		},
	}

	cmd.Flags().StringVar(&opts.host, "host", "", "greenlake host/ip address")
	_ = cmd.MarkFlagRequired("host")
	_ = cmd.MarkFlagRequired("tenantid")

	return cmd
}
