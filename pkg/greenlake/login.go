// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/spf13/cobra"
)

var glLoginData struct {
	host,
	userID,
	secretKey,
	tenantID string
}

// cmdGLLogin represents the green lake login command
var cmdGLLogin = &cobra.Command{
	Use:   "login",
	Short: "Login to greenlake: hpecli greenlake login",
	RunE:  runGLLogin,
}

func init() {
	cmdGLLogin.Flags().StringVar(&glLoginData.host, "host", "", "greenlake host/ip address")
	cmdGLLogin.Flags().StringVarP(&glLoginData.userID, "userid", "u", "", "greenlake userid")
	cmdGLLogin.Flags().StringVarP(&glLoginData.secretKey, "secretkey", "s", "", "greenlake secretkey")
	cmdGLLogin.Flags().StringVarP(&glLoginData.tenantID, "tenantid", "t", "", "greenlake tenantid")
	_ = cmdGLLogin.MarkFlagRequired("host")
	_ = cmdGLLogin.MarkFlagRequired("userid")
	_ = cmdGLLogin.MarkFlagRequired("secretkey")
	_ = cmdGLLogin.MarkFlagRequired("tenantid")
}

func runGLLogin(_ *cobra.Command, _ []string) error {
	if !strings.HasPrefix(glLoginData.host, "http") {
		glLoginData.host = fmt.Sprintf("http://%s", glLoginData.host)
	}

	logger.Debug("Attempting login with user: %v, at: %v", glLoginData.userID, glLoginData.host)

	glc := NewGLClient("client_credentials", glLoginData.userID,
		glLoginData.secretKey, glLoginData.tenantID, glLoginData.host)

	token, err := glc.Login()
	if err != nil {
		logger.Warning("Unable to login with supplied credentials to GreenLake at: %s", glLoginData.host)
		return err
	}

	// change context to current host and save the access token as the API key
	// for subsequent requests
	if err = saveData(glLoginData.host, glLoginData.tenantID, token); err != nil {
		logger.Warning("Successfully logged into GreenLake, but was unable to save the session data")
	} else {
		logger.Debug("Successfully logged into GreenLake: %s", glLoginData.host)
	}

	return nil
}
