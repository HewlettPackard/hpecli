// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package cfm

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newContextCommand() *cobra.Command {
	var host string

	cmd := &cobra.Command{
		Use:   "context",
		Short: "Change context to different CFM host",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runSetContext(host)
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "cfm host/ip address")

	return cmd
}

func runSetContext(host string) error {
	if host != "" {
		return setContext(host)
	}

	// didn't specify host, so just show current context
	ctx, err := getContext()
	if err != nil {
		return err
	}

	logrus.Warningf("Default cfm commands directed to host: %s", ctx)

	return nil
}
