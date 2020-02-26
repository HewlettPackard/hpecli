// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Cmd represents the ilo command
var ovContextCmd = &cobra.Command{
	Use:   "context",
	Short: "Chagne context to different OneView host",
	RunE:  runChangeContext,
}

var ovContextHost struct {
	host string
}

func init() {
	ovContextCmd.Flags().StringVar(&ovContextHost.host, "host", "", "oneview host/ip address")
	_ = ovContextCmd.MarkFlagRequired("host")
}

func runChangeContext(_ *cobra.Command, _ []string) error {
	if !strings.HasPrefix(ovContextHost.host, "http") {
		ovContextHost.host = fmt.Sprintf("https://%s", ovContextHost.host)
	}

	return setContext(ovContextHost.host)
}
