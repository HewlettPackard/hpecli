// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/spf13/cobra"
)

var ovLoginData struct {
	host,
	username,
	password string
}

// ovLoginCmd represents the oneview ovLoginCmd command
var ovLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to OneView: hpecli oneview login",
	RunE:  runOVLogin,
}

func init() {
	ovLoginCmd.Flags().StringVar(&ovLoginData.host, "host", "", "oneview host/ip address")
	ovLoginCmd.Flags().StringVarP(&ovLoginData.username, "username", "u", "", "oneview username")
	ovLoginCmd.Flags().StringVarP(&ovLoginData.password, "password", "p", "", "oneview passowrd")
	_ = ovLoginCmd.MarkFlagRequired("host")
	_ = ovLoginCmd.MarkFlagRequired("username")
	_ = ovLoginCmd.MarkFlagRequired("password")
}

func runOVLogin(_ *cobra.Command, _ []string) error {
	if !strings.HasPrefix(ovLoginData.host, "http") {
		ovLoginData.host = fmt.Sprintf("https://%s", ovLoginData.host)
	}

	logger.Debug("Attempting login with user: %v, at: %v", ovLoginData.username, ovLoginData.host)

	// OneView Login currently doesn't support forced login message acknowledgement - so we roll our own
	token, err := Login(ovLoginData.host, ovLoginData.username, ovLoginData.password)
	if err != nil {
		logger.Warning("Unable to login with supplied credentials to OneView at: %s", ovLoginData.host)
		return err
	}

	// change context to current host and save the session ID as the API key
	// for subsequent requests
	if err = saveContextAndHostData(ovLoginData.host, token); err != nil {
		logger.Warning("Successfully logged into OneView, but was unable to save the session data")
	} else {
		logger.Debug("Successfully logged into OneView: %s", ovLoginData.host)
	}

	return nil
}
