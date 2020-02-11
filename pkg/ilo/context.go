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

var iloContextData struct {
	host string
}

func init() {
	cmdILOContext.Flags().StringVar(&iloContextData.host, "host", "", "ilo host/ip address")
	_ = cmdILOContext.MarkFlagRequired("host")
}

func runSetContext(_ *cobra.Command, _ []string) error {
	if !strings.HasPrefix(iloContextData.host, "http") {
		iloContextData.host = fmt.Sprintf("https://%s", iloContextData.host)
	}

	c := iloContext()

	return c.SetContext(iloContextData.host)
}
