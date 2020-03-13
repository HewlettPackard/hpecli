// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"errors"
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/password"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ovLoginOptions struct {
	host,
	username,
	password string
	passwordStdin bool
}

func newLoginCommand() *cobra.Command {
	var opts ovLoginOptions

	var cmd = &cobra.Command{
		Use:   "login",
		Short: "Login to OneView: hpecli oneview login",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return validateArgs(&opts)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return runLogin(&opts)
		},
	}

	cmd.Flags().StringVar(&opts.host, "host", "", "oneview host/ip address")
	cmd.Flags().StringVarP(&opts.username, "username", "u", "", "oneview username")
	cmd.Flags().StringVarP(&opts.password, "password", "p", "", "oneview password")
	cmd.Flags().BoolVarP(&opts.passwordStdin, "password-stdin", "", false, "read password from stdin")
	_ = cmd.MarkFlagRequired("host")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}

func validateArgs(opts *ovLoginOptions) error {
	if opts.host != "" && !strings.HasPrefix(opts.host, "http") {
		opts.host = fmt.Sprintf("https://%s", opts.host)
	}

	if opts.password != "" && opts.passwordStdin {
		return errors.New("--password and --password-stdin are mutually exclusive")
	}

	return nil
}

func runLogin(opts *ovLoginOptions) error {
	if err := password.Read(&opts.password, opts.passwordStdin, "oneview password: "); err != nil {
		return err
	}

	logrus.Debugf("Attempting login with user: %v, at: %v", opts.username, opts.host)

	// OneView Login currently doesn't support forced login message acknowledgement - so we roll our own
	token, err := login(opts.host, opts.username, opts.password)
	if err != nil {
		logrus.Warningf("Unable to login with supplied credentials to OneView at: %s", opts.host)
		return err
	}

	// change context to current host and save the session ID as the API key
	// for subsequent requests
	if err = saveContextAndHostData(opts.host, token); err != nil {
		logrus.Warning("Successfully logged into OneView, but was unable to save the session data")
	} else {
		logrus.Warningf("Successfully logged into OneView: %s", opts.host)
	}

	return nil
}
