// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/HewlettPackard/hpecli/internal/platform/password"
	"github.com/spf13/cobra"
)

type cvLoginOptions struct {
	host     string
	username string
	password string
}

func newLoginCommand() *cobra.Command {
	var opts cvLoginOptions

	// Cmd represents the cloudvolume get command
	var cmd = &cobra.Command{
		Use:   "login",
		Short: "Login to HPE Nimble Cloud Volumes: hpecli cloudvolume login",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(&opts)
		},
	}

	cmd.Flags().StringVar(&opts.host, "host", "", "Cloud Volumes portal hostname/ip")
	cmd.Flags().StringVarP(&opts.username, "username", "u", "", "cloudvolume username")
	cmd.Flags().StringVarP(&opts.password, "password", "p", "", "cloudvolume passowrd")
	_ = cmd.MarkFlagRequired("host")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}

func runLogin(opts *cvLoginOptions) error {
	log.Logger.Debug("cloudvolumes/login called")

	if !strings.HasPrefix(opts.host, "http") {
		opts.host = fmt.Sprintf("https://%s", opts.host)
	}

	if opts.password == "" {
		p, err := password.ReadFromConsole("cloudvolumes password: ")
		if err != nil {
			log.Logger.Error("Unable to read password from console!")
			return err
		}

		opts.password = p
	}

	log.Logger.Debugf("Attempting login with user: %v, at: %v", opts.username, opts.host)

	cvc := newCVClient(opts.host, opts.username, opts.password)

	token, err := cvc.Login()
	if err != nil {
		log.Logger.Warningf("Unable to login with supplied credentials to CloudVolume at: %s", opts.host)
		return err
	}

	// change context to current host and save the token as the API key
	// for subsequent requests
	if err := saveData(opts.host, token); err != nil {
		log.Logger.Warning("Successfully logged into CloudVolumes, but was unable to save the session data")
		log.Logger.Debugf("%+v", err)
	} else {
		log.Logger.Warningf("Successfully logged into CloudVolumes: %s", opts.host)
	}

	return nil
}
