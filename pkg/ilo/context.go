// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Cmd represents the ilo command
var cmdILOContext = &cobra.Command{
	Use:   "context",
	Short: "Chagne context to different ilo host",
	RunE:  runSetContext,
}

var iloContextHost struct {
	host string
}

func init() {
	cmdILOContext.Flags().StringVar(&iloContextHost.host, "host", "", "ilo host/ip address")
	_ = cmdILOContext.MarkFlagRequired("host")
}

func runSetContext(_ *cobra.Command, _ []string) error {
	if !strings.HasPrefix(iloContextHost.host, "http") {
		iloContextHost.host = fmt.Sprintf("https://%s", iloContextHost.host)
	}

	return changeContext(iloContextHost.host)
}
