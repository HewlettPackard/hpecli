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
	tenantID,
	grantType string
}

func init() {
	glLoginCmd.Flags().StringVar(&glLoginData.host, "host", "", "greenlake host/ip address")
	glLoginCmd.Flags().StringVarP(&glLoginData.userID, "userid", "u", "", "greenlake userid")
	glLoginCmd.Flags().StringVarP(&glLoginData.secretKey, "secretkey", "s", "", "greenlake secretkey")
	glLoginCmd.Flags().StringVarP(&glLoginData.tenantID, "tenantid", "t", "", "greenlake tenantid")
	_ = glLoginCmd.MarkFlagRequired("host")
	_ = glLoginCmd.MarkFlagRequired("userid")
	_ = glLoginCmd.MarkFlagRequired("secretkey")
	_ = glLoginCmd.MarkFlagRequired("tenantid")
}

// getCmd represents the get command
var glLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to greenlake: hpecli greenlake login",
	RunE:  runGLLogin,
}

func runGLLogin(_ *cobra.Command, _ []string) error {

	logger.Info("greenlake/login called")
	//cmd.SilenceUsage = true
	//cmd.SilenceErrors = true

	if !strings.HasPrefix(glLoginData.host, "http") {
		glLoginData.host = fmt.Sprintf("http://%s", glLoginData.host)
	}

	logger.Debug("Attempting login with user: %v, at: %v", glLoginData.userID, glLoginData.host)

	glc := NewGreenLakeClient("client_credentials", glLoginData.userID, glLoginData.secretKey, glLoginData.tenantID, glLoginData.host)

	s, err := glc.GetToken()
	if err != nil {
		logger.Warning("Unable to login with supplied credentials to GreenLake at: %s", glLoginData.host)
		return err
	}

	glc.APIKey = s.AccessToken
	glc.TenantID = glLoginData.tenantID
	println("Access Token", s.AccessToken)

	// change context to current host and save the session ID as the API key
	// for subsequent requests
	if err = setTokenTentanID(glLoginData.host, glLoginData.tenantID, s.AccessToken); err != nil {
		logger.Warning("Successfully logged into GreenLake, but was unable to save the session data")
	} else {
		println("Successfully logged into GreenLake: ", glLoginData.host)
		logger.Debug("Successfully logged into GreenLake: %s", glLoginData.host)
	}

	return nil
}
