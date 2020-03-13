// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"errors"
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/password"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type cvLoginOptions struct {
	host          string
	username      string
	password      string
	passwordStdin bool
}

func newLoginCommand() *cobra.Command {
	var opts cvLoginOptions

	// Cmd represents the cloudvolume get command
	cmd := &cobra.Command{
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return validateArgs(&opts)
		},
		Use:   "login",
		Short: "Login to HPE Nimble Cloud Volumes: hpecli cloudvolumes login",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runLogin(&opts)
		},
	}

	cmd.Flags().StringVar(&opts.host, "host", "", "Cloud Volumes portal hostname/ip")
	cmd.Flags().StringVarP(&opts.username, "username", "u", "", "cloudvolume username")
	cmd.Flags().StringVarP(&opts.password, "password", "p", "", "cloudvolume passowrd")
	cmd.Flags().BoolVarP(&opts.passwordStdin, "password-stdin", "", false, "read password from stdin")
	_ = cmd.MarkFlagRequired("host")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}

func validateArgs(opts *cvLoginOptions) error {
	if opts.host != "" && !strings.HasPrefix(opts.host, "http") {
		opts.host = fmt.Sprintf("https://%s", opts.host)
	}

	if opts.password != "" && opts.passwordStdin {
		return errors.New("--password and --password-stdin are mutually exclusive")
	}

	return nil
}

func runLogin(opts *cvLoginOptions) error {
	logrus.Debug("cloudvolumes/login called")

	if err := password.Read(&opts.password, opts.passwordStdin, "cloudvolumes password: "); err != nil {
		return err
	}

	logrus.Debugf("Attempting login with user: %v, at: %v", opts.username, opts.host)

	cvc := newCVClient(opts.host, opts.username, opts.password)

	token, err := cvc.login()
	if err != nil {
		logrus.Warningf("Unable to login with supplied credentials to CloudVolume at: %s", opts.host)
		return err
	}

	// change context to current host and save the token as the API key
	// for subsequent requests
	if err := saveData(opts.host, token); err != nil {
		logrus.Warning("Successfully logged into CloudVolumes, but was unable to save the session data")
		logrus.Debugf("%+v", err)
	} else {
		logrus.Warningf("Successfully logged into CloudVolumes: %s", opts.host)
	}

	return nil
}
