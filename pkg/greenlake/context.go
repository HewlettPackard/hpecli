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

var glContextHostTenant struct {
	host     string
	tenantID string
}

func init() {
	cmdGLContext.Flags().StringVar(&glContextHostTenant.host, "host", "", "greenlake host/ip address")
	cmdGLContext.Flags().StringVar(&glContextHostTenant.tenantID, "tenantid", "t", "greenlake tenantid")
	_ = cmdGLContext.MarkFlagRequired("host")
	_ = cmdGLContext.MarkFlagRequired("tenantid")
}

func runSetContext(_ *cobra.Command, _ []string) error {
	if !strings.HasPrefix(glContextHostTenant.host, "http") {
		glContextHostTenant.host = fmt.Sprintf("https://%s", glContextHostTenant.host)
	}

	return changeContext(glContextHostTenant.host)
}
