// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newContextCommand() *cobra.Command {
	var host string

	cmd := &cobra.Command{
		Use:   "context",
		Short: "Chagne context to different OneView host",
		RunE: func(_ *cobra.Command, _ []string) error {
			if host != "" && !strings.HasPrefix(host, "http") {
				host = fmt.Sprintf("https://%s", host)
			}

			return runSetContext(host)
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "oneview host/ip address")

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

	logrus.Warningf("Default oneview commands directed to host: %s", ctx)

	return nil
}
