// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// cmdGLContext represents the green lake command
var cmdGLContext = &cobra.Command{
	Use:   "context",
	Short: "Change context to different GreenLake host",
	RunE:  runSetContext,
}

var glContextHost struct {
	host string
}

func init() {
	cmdGLContext.Flags().StringVar(&glContextHost.host, "host", "", "greenlake host/ip address")
	_ = cmdGLContext.MarkFlagRequired("host")
	_ = cmdGLContext.MarkFlagRequired("tenantid")
}

func runSetContext(_ *cobra.Command, _ []string) error {
	if !strings.HasPrefix(glContextHost.host, "http") {
		glContextHost.host = fmt.Sprintf("https://%s", glContextHost.host)
	}

	return changeContext(glContextHost.host)
}
