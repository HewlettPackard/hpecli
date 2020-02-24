// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/spf13/cobra"
)

var cvLoginData struct {
	host,
	username,
	password,
	token string
}

// Cmd represents the cloudvolume get command
var cmdLogin = &cobra.Command{
	Use:   "login",
	Short: "Login to HPE Nimble Cloud Volumes: hpecli cloudvolume login",
	RunE:  runCVLogin,
}

func init() {
	cmdLogin.Flags().StringVar(&cvLoginData.host, "host", "", "Cloud Volumes portal hostname/ip")
	cmdLogin.Flags().StringVarP(&cvLoginData.username, "username", "u", "", "cloudvolume username")
	cmdLogin.Flags().StringVarP(&cvLoginData.password, "password", "p", "", "cloudvolume passowrd")
	_ = cmdLogin.MarkFlagRequired("host")
	_ = cmdLogin.MarkFlagRequired("username")
	_ = cmdLogin.MarkFlagRequired("password")
}

func runCVLogin(_ *cobra.Command, _ []string) error {
	logger.Debug("cloudvolumes/login called")

	if !strings.HasPrefix(cvLoginData.host, "http") {
		cvLoginData.host = fmt.Sprintf("https://%s", cvLoginData.host)
	}

	logger.Debug("Attempting login with user: %v, at: %v", cvLoginData.username, cvLoginData.host)

	cvc := NewCVClient(cvLoginData.host, cvLoginData.username, cvLoginData.password)

	token, err := cvc.Login()
	if err != nil {
		logger.Warning("Unable to login with supplied credentials to CloudVolume at: %s", cvLoginData.host)
		return err
	}

	// change context to current host and save the token as the API key
	// for subsequent requests
	if err := saveData(cvLoginData.host, token); err != nil {
		logger.Warning("Successfully logged into CloudVolumes, but was unable to save the session data")
		logger.Debug("%+v", err)
	} else {
		logger.Debug("Successfully logged into CloudVolumes: %s", cvLoginData.host)
	}

	return nil
}
