// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newContextCommand() *cobra.Command {
	var host string

	var cmd = &cobra.Command{
		Use:   "context",
		Short: "Chagne context to different ilo host",
		RunE: func(cmd *cobra.Command, args []string) error {
			if host != "" && !strings.HasPrefix(host, "http") {
				host = fmt.Sprintf("https://%s", host)
			}
			return runSetContext(host)
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "ilo host/ip address")

	return cmd
}

func runSetContext(host string) error {
	// didn't specify host, so just show current context
	if host == "" {
		ctx, err := getContext()
		if err != nil {
			return err
		}

		logrus.Warningf("Default ilo commands directed to host: %s", ctx)

		return nil
	}

	return setContext(host)
}
