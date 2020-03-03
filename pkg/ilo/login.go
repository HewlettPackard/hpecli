// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/HewlettPackard/hpecli/internal/platform/password"
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
}

func runILOLogin(_ *cobra.Command, _ []string) error {
	if !strings.HasPrefix(iloLoginData.host, "http") {
		iloLoginData.host = fmt.Sprintf("https://%s", iloLoginData.host)
	}

	if iloLoginData.password == "" {
		p, err := password.ReadFromConsole("ilo password: ")
		if err != nil {
			log.Logger.Errorln("\nUnable to read password from console!")
			return err
		}

		iloLoginData.password = p
	}

	log.Logger.Debugf("Attempting login with user: %v, at: %v", iloLoginData.username, iloLoginData.host)

	cl := NewILOClient(iloLoginData.host, iloLoginData.username, iloLoginData.password)

	sd, err := cl.login()
	if err != nil {
		log.Logger.Warningf("Unable to login with supplied credentials to ilo at: %s", iloLoginData.host)
		return err
	}

	// change context to current host and save the session ID as the API key
	// for subsequent requests
	if err = saveContextAndSessionData(sd); err != nil {
		log.Logger.Warning("Successfully logged into ilo, but was unable to save the session data")
	} else {
		log.Logger.Debugf("Successfully logged into ilo: %s", iloLoginData.host)
	}

	return nil
}
