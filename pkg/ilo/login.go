// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/spf13/cobra"
)

var iloLoginData struct {
	host,
	username,
	password string
}

// cmdIloLogin represents the get command
var cmdILOLogin = &cobra.Command{
	Use:   "login",
	Short: "Login to iLO: hpecli ilo login",
	RunE:  runILOLogin,
}

func init() {
	cmdILOLogin.Flags().StringVar(&iloLoginData.host, "host", "", "ilo ip address")
	cmdILOLogin.Flags().StringVarP(&iloLoginData.username, "username", "u", "", "ilo username")
	cmdILOLogin.Flags().StringVarP(&iloLoginData.password, "password", "p", "", "ilo passowrd")
	_ = cmdILOLogin.MarkFlagRequired("host")
	_ = cmdILOLogin.MarkFlagRequired("username")
	_ = cmdILOLogin.MarkFlagRequired("password")
}

func runILOLogin(_ *cobra.Command, _ []string) error {
	if !strings.HasPrefix(iloLoginData.host, "http") {
		iloLoginData.host = fmt.Sprintf("https://%s", iloLoginData.host)
	}

	logger.Debug("Attempting login with user: %v, at: %v", iloLoginData.username, iloLoginData.host)

	cl := NewILOClient(iloLoginData.host, iloLoginData.username, iloLoginData.password)

	token, err := cl.Login()
	if err != nil {
		logger.Warning("Unable to login with supplied credentials to ilo at: %s", iloLoginData.host)
		return err
	}

	// change context to current host and save the session ID as the API key
	// for subsequent requests
	if err = storeContext(iloLoginData.host, token); err != nil {
		logger.Warning("Successfully logged into ilo, but was unable to save the session data")
	} else {
		logger.Debug("Successfully logged into ilo: %s", iloLoginData.host)
	}

	return nil
}
