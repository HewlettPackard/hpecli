// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"errors"
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/HewlettPackard/hpecli/internal/platform/password"
	"github.com/spf13/cobra"
)

var ovLoginData struct {
	host,
	username,
	password string
	passwordStdin bool
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
	ovLoginCmd.Flags().StringVarP(&ovLoginData.password, "password", "p", "", "oneview password")
	ovLoginCmd.Flags().BoolVarP(&ovLoginData.passwordStdin, "password-stdin", "", false, "read password from stdin")
	_ = ovLoginCmd.MarkFlagRequired("host")
	_ = ovLoginCmd.MarkFlagRequired("username")
}

func runOVLogin(_ *cobra.Command, _ []string) error {
	if !strings.HasPrefix(ovLoginData.host, "http") {
		ovLoginData.host = fmt.Sprintf("https://%s", ovLoginData.host)
	}

	if err := handlePasswordOptions(); err != nil {
		return err
	}

	log.Logger.Debugf("Attempting login with user: %v, at: %v", ovLoginData.username, ovLoginData.host)

	// OneView Login currently doesn't support forced login message acknowledgement - so we roll our own
	token, err := Login(ovLoginData.host, ovLoginData.username, ovLoginData.password)
	if err != nil {
		log.Logger.Warningf("Unable to login with supplied credentials to OneView at: %s", ovLoginData.host)
		return err
	}

	// change context to current host and save the session ID as the API key
	// for subsequent requests
	if err = saveContextAndHostData(ovLoginData.host, token); err != nil {
		log.Logger.Warning("Successfully logged into OneView, but was unable to save the session data")
	} else {
		log.Logger.Warningf("Successfully logged into OneView: %s", ovLoginData.host)
	}

	return nil
}

func handlePasswordOptions() error {
	if ovLoginData.password != "" {
		if ovLoginData.passwordStdin {
			return errors.New("--password and --password-stdin are mutually exclusive")
		}
		// if the password was set .. we don't need to get it from somewhere else
		return nil
	}

	// asked to read from stdin
	if ovLoginData.passwordStdin {
		pswd, err := password.ReadFromStdIn()
		if err != nil {
			return err
		}

		ovLoginData.password = pswd

		return nil
	}

	// password wasn't specified so we need to prompt them for it
	pswd, err := password.ReadFromConsole("oneview password: ")
	if err != nil {
		return err
	}

	ovLoginData.password = pswd

	return nil
}
