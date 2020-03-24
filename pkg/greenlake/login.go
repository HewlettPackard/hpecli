// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"errors"
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/password"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type glLoginOptions struct {
	host,
	userID,
	secretKey,
	tenantID string
	secretKeyStdin bool
}

func newLoginCommand() *cobra.Command {
	var opts glLoginOptions

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to greenlake: hpecli greenlake login",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return validateArgs(&opts)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return runLogin(&opts)
		},
	}

	cmd.Flags().StringVar(&opts.host, "host", "", "greenlake host/ip address")
	cmd.Flags().StringVarP(&opts.secretKey, "secretkey", "s", "", "greenlake secretkey")
	cmd.Flags().BoolVarP(&opts.secretKeyStdin, "secretkey-stdin", "", false, "read secretkey from stdin")
	cmd.Flags().StringVarP(&opts.tenantID, "tenantid", "t", "", "greenlake tenantid")
	cmd.Flags().StringVarP(&opts.userID, "userid", "u", "", "greenlake userid")
	_ = cmd.MarkFlagRequired("host")
	_ = cmd.MarkFlagRequired("tenantid")
	_ = cmd.MarkFlagRequired("userid")

	return cmd
}

func validateArgs(opts *glLoginOptions) error {
	if opts.host != "" && !strings.HasPrefix(opts.host, "http") {
		opts.host = fmt.Sprintf("http://%s", opts.host)
	}

	if opts.secretKey != "" && opts.secretKeyStdin {
		return errors.New("--secretKey and --secretkey-stdin are mutually exclusive")
	}

	return nil
}

func runLogin(opts *glLoginOptions) error {
	if err := password.Read(&opts.secretKey, opts.secretKeyStdin, "greenlake secretKey: "); err != nil {
		return err
	}

	logrus.Debugf("Attempting login with user: %v, at: %v", opts.userID, opts.host)

	glc := newGLClient("client_credentials", opts.userID, opts.secretKey, opts.tenantID, opts.host)

	sd, err := glc.login()
	if err != nil {
		logrus.Warningf("Unable to login with supplied credentials to GreenLake at: %s", opts.host)
		return err
	}

	// change context to current host and save the access token as the API key
	// for subsequent requests
	if err = saveContextAndSessionData(sd); err != nil {
		logrus.Debug("Successfully logged into GreenLake, but was unable to save the session data")
	} else {
		logrus.Infof("Successfully logged into GreenLake: %s", opts.host)
	}

	return nil
}
