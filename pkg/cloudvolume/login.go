// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/HewlettPackard/hpecli/internal/platform/password"
	"github.com/spf13/cobra"
)

var cvLoginData struct {
	host,
	username,
	password string
}

// Cmd represents the cloudvolume get command
var cmdLogin = &cobra.Command{
	Use:   "login",
	Short: "Login to HPE Nimble Cloud Volumes: hpecli cloudvolume login",
	RunE:  runCVLogin,
}

func init() {
	cmdLogin.Flags().StringVar(&cvLoginData.host, "host", "", "Cloud Volumes portal hostname/ip")
	cmdLogin.Flags().StringVarP(&cvLoginData.username, "username", "u", "", "cloudvolume username")
	cmdLogin.Flags().StringVarP(&cvLoginData.password, "password", "p", "", "cloudvolume passowrd")
	_ = cmdLogin.MarkFlagRequired("host")
	_ = cmdLogin.MarkFlagRequired("username")
}

func runCVLogin(_ *cobra.Command, _ []string) error {
	log.Logger.Debug("cloudvolumes/login called")

	if !strings.HasPrefix(cvLoginData.host, "http") {
		cvLoginData.host = fmt.Sprintf("https://%s", cvLoginData.host)
	}

	if cvLoginData.password == "" {
		p, err := password.ReadFromConsole("cloudvolumes password: ")
		if err != nil {
			log.Logger.Errorln("\nUnable to read password from console!")
			return err
		}

		cvLoginData.password = p
	}

	log.Logger.Debugf("Attempting login with user: %v, at: %v", cvLoginData.username, cvLoginData.host)

	cvc := NewCVClient(cvLoginData.host, cvLoginData.username, cvLoginData.password)

	token, err := cvc.Login()
	if err != nil {
		log.Logger.Warningf("Unable to login with supplied credentials to CloudVolume at: %s", cvLoginData.host)
		return err
	}

	// change context to current host and save the token as the API key
	// for subsequent requests
	if err := saveData(cvLoginData.host, token); err != nil {
		log.Logger.Warning("Successfully logged into CloudVolumes, but was unable to save the session data")
		log.Logger.Debugf("%+v", err)
	} else {
		log.Logger.Debugf("Successfully logged into CloudVolumes: %s", cvLoginData.host)
	}

	return nil
}
