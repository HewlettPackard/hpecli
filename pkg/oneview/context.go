// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
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
}

func runChangeContext(_ *cobra.Command, _ []string) error {
	// didn't specify host, so just show current context
	if ovContextHost.host == "" {
		ctx, err := getContext()
		if err != nil {
			return err
		}

		log.Logger.Warningf("Default oneview commands directed to host: %s", ctx)

		return nil
	}

	if !strings.HasPrefix(ovContextHost.host, "http") {
		ovContextHost.host = fmt.Sprintf("https://%s", ovContextHost.host)
	}

	return setContext(ovContextHost.host)
}
