// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newStatusCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "status",
		Short: "Check whether google analytics is on or off",
		RunE: func(cmd *cobra.Command, _ []string) error {
			checkAnalytics := CheckGoogleAnalytics()
			eventCategory := "analytics"
			eventAction := cmd.Name()
			err := checkStatus(cmd, checkAnalytics, eventCategory, eventAction)
			return err
		},
	}

	return cmd
}

func checkStatus(cmd *cobra.Command, checkAnalytics bool, eventCategory, eventAction string) error {
	analyticsClient := NewAnalyticsClient("1", "event", eventCategory,
		eventAction, "200", "", "hpecli/0.0.1", "0.0.1", "hpecli")
	_, err := analyticsClient.TrackEvent()

	if err != nil {
		logrus.Warningf("Unable to send the analytics info with supplied event details")
		return err
	}

	if checkAnalytics {
		logrus.Info(" Google Analytics is turned ON \n please run \"hpe analytics off\", if you want to turn it off")
	} else {
		logrus.Info(" Google Analytics is turned OFF \n please run \"hpe analytics on\", if you want to turn it on")
	}

	return nil
}
