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
	RunE:  runSetContext,
}

var ovContextData struct {
	host string
}

func init() {
	ovContextCmd.Flags().StringVar(&ovContextData.host, "host", "", "oneview host/ip address")
	_ = ovContextCmd.MarkFlagRequired("host")
}

func runSetContext(_ *cobra.Command, _ []string) error {
	if !strings.HasPrefix(ovContextData.host, "http") {
		ovContextData.host = fmt.Sprintf("https://%s", ovContextData.host)
	}

	c := ovContext()

	return c.SetContext(ovContextData.host)
}
