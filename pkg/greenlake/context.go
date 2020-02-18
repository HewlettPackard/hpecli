// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Cmd represents the ilo command
var glContextCmd = &cobra.Command{
	Use:   "context",
	Short: "Chagne context to different GreenLake host",
	RunE:  runChangeContext,
}

var glContextData struct {
	host     string
	tenantID string
}

func init() {
	glContextCmd.Flags().StringVar(&glContextData.host, "host", "", "greenlake host/ip address")
	glLoginCmd.Flags().StringVar(&glContextData.tenantID, "tenantid", "t", "", "greenlake tenantid")
	_ = glContextCmd.MarkFlagRequired("host")
	_ = glContextCmd.MarkFlagRequired("tenantid")
}

func runChangeContext(_ *cobra.Command, _ []string) error {
	if !strings.HasPrefix(glContextData.host, "http") {
		glContextData.host = fmt.Sprintf("https://%s", glContextData.host)
	}

	c := glContext()

	return changeContext(glContextData.host)
}
