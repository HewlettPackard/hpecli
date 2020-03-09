// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newUsersCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "users",
		Short: "Get from Users: hpecli greenlake get users",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runUsers()
		},
	}

	return cmd
}

func runUsers() error {
	logrus.Debug("Beginning runUsers")

	sd, err := defaultSessionData()
	if err != nil {
		logrus.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for HPE GreenLake.  " +
			"Please login to GreenLake using: hpecli greenlake login")
	}

	logrus.Debugf("Attempting get green lake users at: %v", sd.Host)

	glc := newGLClientFromAPIKey(sd.Host, sd.TenantID, sd.Token)

	jsonResult, err := glc.users()
	if err != nil {
		return err
	}

	logrus.Infof("%s", jsonResult)

	return nil
}
