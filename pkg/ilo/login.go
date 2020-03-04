// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"errors"
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
	passwordStdin bool
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
	cmdILOLogin.Flags().StringVarP(&iloLoginData.password, "password", "p", "", "ilo password")
	cmdILOLogin.Flags().BoolVarP(&iloLoginData.passwordStdin, "password-stdin", "", false, "read password from stdin")
	_ = cmdILOLogin.MarkFlagRequired("host")
	_ = cmdILOLogin.MarkFlagRequired("username")
}

func runILOLogin(_ *cobra.Command, _ []string) error {
	if !strings.HasPrefix(iloLoginData.host, "http") {
		iloLoginData.host = fmt.Sprintf("https://%s", iloLoginData.host)
	}

	if err := handlePasswordOptions(); err != nil {
		return err
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
		log.Logger.Warningf("Successfully logged into ilo: %s", iloLoginData.host)
	}

	return nil
}

func handlePasswordOptions() error {
	if iloLoginData.password != "" {
		if iloLoginData.passwordStdin {
			return errors.New("--password and --password-stdin are mutually exclusive")
		}
		// if the password was set .. we don't need to get it from somewhere else
		return nil
	}

	// asked to read from stdin
	if iloLoginData.passwordStdin {
		pswd, err := password.ReadFromStdIn()
		if err != nil {
			return err
		}

		iloLoginData.password = pswd

		return nil
	}

	// password wasn't specified so we need to prompt them for it
	pswd, err := password.ReadFromConsole("ilo password: ")
	if err != nil {
		return err
	}

	iloLoginData.password = pswd

	return nil
}
