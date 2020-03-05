// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newContextCommand() *cobra.Command {
	var host string

	var cmd = &cobra.Command{
		Use:   "context",
		Short: "Chagne context to different ilo host",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetContext(&host)
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "ilo host/ip address")
	_ = cmd.MarkFlagRequired("host")

	return cmd
}

func runSetContext(host *string) error {
	if !strings.HasPrefix(*host, "http") {
		*host = fmt.Sprintf("https://%s", *host)
	}

	return setContext(*host)
}
