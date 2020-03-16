// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newServersCommand() *cobra.Command {
	var serverName string

	var cmd = &cobra.Command{
		Use:   "servers",
		Short: "Get servers from OneView: hpecli oneview get servers",
		RunE: func(_ *cobra.Command, _ []string) error {
			return getServerHardware(serverName)
		},
	}

	cmd.Flags().StringVar(&serverName, "name", "", "name of the server to retrieve")

	return cmd
}

func getServerHardware(serverName string) error {
	host, token, err := hostAndToken()
	if err != nil {
		logrus.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for OneView.  " +
			"Please login to OneView using: hpecli oneview login")
	}

	ovc := newOVClientFromAPIKey(host, token)

	logrus.Warningf("Using OneView: %s", host)

	var sh interface{}
	if serverName != "" {
		sh, err = ovc.GetServerHardwareByName(serverName)
	} else {
		sh, err = ovc.GetServerHardwareList(nil, "", "", "", "")
	}

	if err != nil {
		logrus.Warningf("Unable to login with supplied credentials to OneView at: %s", host)
		return err
	}

	out, err := json.MarshalIndent(sh, "", "  ")
	if err != nil {
		logrus.Warning("Unable to output data as JSON.  Please try the command again.")
	}

	logrus.Infof("%s", out)

	return nil
}
