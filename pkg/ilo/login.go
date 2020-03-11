// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/password"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type iloLoginOptions struct {
	host,
	username,
	password string
	passwordStdin bool
}

func newLoginCommand() *cobra.Command {
	var opts iloLoginOptions

	var cmd = &cobra.Command{
		Use:   "login",
		Short: "Login to iLO: hpecli ilo login",
		RunE: func(_ *cobra.Command, _ []string) error {
			if !strings.HasPrefix(opts.host, "http") {
				opts.host = fmt.Sprintf("https://%s", opts.host)
			}

			return runLogin(&opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.passwordStdin, "password-stdin", "", false, "read password from stdin")
	cmd.Flags().StringVar(&opts.host, "host", "", "ilo ip address")
	cmd.Flags().StringVarP(&opts.username, "username", "u", "", "ilo username")
	cmd.Flags().StringVarP(&opts.password, "password", "p", "", "ilo password")
	_ = cmd.MarkFlagRequired("host")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}

func runLogin(opts *iloLoginOptions) error {
	if err := handlePasswordOptions(opts); err != nil {
		return err
	}

	logrus.Debugf("Attempting login with user: %v, at: %v", opts.username, opts.host)

	cl := newILOClient(opts.host, opts.username, opts.password)

	sd, err := cl.login()
	if err != nil {
		logrus.Warningf("Unable to login with supplied credentials to ilo at: %s", opts.host)
		return err
	}

	// change context to current host and save the session ID as the API key
	// for subsequent requests
	if err = saveContextAndSessionData(sd); err != nil {
		logrus.Warning("Successfully logged into ilo, but was unable to save the session data")
	} else {
		logrus.Warningf("Successfully logged into ilo: %s", opts.host)
	}

	return nil
}

func handlePasswordOptions(opts *iloLoginOptions) error {
	if opts.password != "" {
		if opts.passwordStdin {
			return errors.New("--password and --password-stdin are mutually exclusive")
		}
		// if the password was set .. we don't need to get it from somewhere else
		return nil
	}

	// asked to read from stdin
	if opts.passwordStdin {
		pswd, err := password.ReadFromStdIn()
		if err != nil {
			return err
		}

		opts.password = pswd

		return nil
	}

	// password wasn't specified so we need to prompt them for it
	pswd, err := password.ReadFromConsole("ilo password: ")
	if err != nil {
		return err
	}

	opts.password = pswd

	return nil
}
