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
		Short: "Login to iLO",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return validateArgs(&opts)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return runLogin(&opts)
		},
	}

	cmd.Flags().StringVar(&opts.host, "host", "", "iLO ip address")
	cmd.Flags().StringVarP(&opts.username, "username", "u", "", "iLO username")
	cmd.Flags().StringVarP(&opts.password, "password", "p", "", "iLO password")
	cmd.Flags().BoolVarP(&opts.passwordStdin, "password-stdin", "", false, "read password from stdin")
	_ = cmd.MarkFlagRequired("host")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}

func validateArgs(opts *iloLoginOptions) error {
	if opts.host != "" && !strings.HasPrefix(opts.host, "http") {
		opts.host = fmt.Sprintf("https://%s", opts.host)
	}

	if opts.password != "" {
		logrus.Warning("WARNING! Using --password via the CLI is insecure. Use --password-stdin.")

		if opts.passwordStdin {
			return errors.New("--password and --password-stdin are mutually exclusive")
		}
	}

	return nil
}

func runLogin(opts *iloLoginOptions) error {
	if err := password.Read(&opts.password, opts.passwordStdin, "iLO password: "); err != nil {
		return err
	}

	logrus.Debugf("Attempting login with user: %v, at: %v", opts.username, opts.host)

	cl := newILOClient(opts.host, opts.username, opts.password)

	sd, err := cl.login()
	if err != nil {
		logrus.Warningf("Unable to login with supplied credentials to iLO at: %s", opts.host)
		return err
	}

	// change context to current host and save the session ID as the API key
	// for subsequent requests
	if err = saveContextAndSessionData(sd); err != nil {
		logrus.Warning("Successfully logged into iLO, but was unable to save the session data")
	} else {
		logrus.Infof("Successfully logged into iLO: %s", opts.host)
	}

	return nil
}
