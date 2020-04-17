// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newOffCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "off",
		Short: "Turn off Google Analytics",
		RunE: func(cmd *cobra.Command, _ []string) error {
			disable, disableerr := disableGoogleAnalytics()
			if disableerr != nil {
				logrus.Debugf("Error disabling Google analytics and the error is %s", disableerr)
			}
			eventCategory := "analytics"
			eventAction := cmd.Name()
			err := offAnalytics(cmd, disable, eventCategory, eventAction)
			if err != nil {
				logrus.Debugf("Unable to turn off the Google Analytics and the error is %s", err)
				return err
			}
			return nil
		},
	}

	return cmd
}

func offAnalytics(cmd *cobra.Command, disable bool, eventCategory, eventAction string) error {
	if disable {
		analyticsClient := NewAnalyticsClient("1", "event", eventCategory,
			eventAction, "200", "", "hpecli/0.0.1", "0.0.1", "hpecli")
		err := analyticsClient.TrackEvent()

		if err != nil {
			logrus.Warningf("Unable to send the GA info with supplied event details")
			return err
		}

		logrus.Info(" Google Analytics is turned OFF \n please run \"hpe analytics on\", if you want to turn it on")
	}

	return nil
}
