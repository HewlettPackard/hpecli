// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"fmt"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/spf13/cobra"
)

// cmdGetUsers represents the green lake users command
var cmdGetUsers = &cobra.Command{
	Use:   "users",
	Short: "Get from Users: hpecli greenlake get users",
	RunE:  runGLGetUsers,
}

func runGLGetUsers(_ *cobra.Command, _ []string) error {
	log.Logger.Debug("Beginning runGLGetUsers")

	sd, err := defaultSessionData()
	if err != nil {
		log.Logger.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for HPE GreenLake." +
			"Please login to GreenLake using: hpecli greenlake login")
	}

	log.Logger.Debugf("Attempting get green lake users at: %v", sd.Host)

	glc := NewGLClientFromAPIKey(sd.Host, sd.TenantID, sd.Token)

	jsonResult, err := glc.GetUsers()
	if err != nil {
		return err
	}

	log.Logger.Infof("%s", jsonResult)

	return nil
}
