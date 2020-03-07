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
		RunE: func(_ *cobra.Command, _ []string) error {
			if !strings.HasPrefix(opts.host, "http") {
				opts.host = fmt.Sprintf("https://%s", opts.host)
			}
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

func runLogin(opts *ovLoginOptions) error {
	if err := handlePasswordOptions(opts); err != nil {
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

func handlePasswordOptions(opts *ovLoginOptions) error {
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
	pswd, err := password.ReadFromConsole("oneview password: ")
	if err != nil {
		return err
	}

	opts.password = pswd

	return nil
}
