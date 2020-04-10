// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const analytics = "analytics"

func newOnCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "on",
		Short: "Turn on Google Analytics",
		RunE: func(cmd *cobra.Command, _ []string) error {
			enable := enableGoogleAnalytics()
			eventCategory := analytics
			eventAction := cmd.Name()
			err := onAnalytics(cmd, enable, eventCategory, eventAction)
			return err
		},
	}

	return cmd
}

func onAnalytics(cmd *cobra.Command, enable bool, eventCategory, eventAction string) error {
	if enable {
		analyticsClient := NewAnalyticsClient("1", "event", eventCategory,
			eventAction, "200", "", "hpecli/0.0.1", "0.0.1", "hpecli")
		_, err := analyticsClient.TrackEvent()

		if err != nil {
			logrus.Warning("Unable to send the GA info with supplied event details")
			return err
		}

		logrus.Info(" Google Analytics is turned ON \n please run \"hpe analytics off\", if you want to turn it off")
	}

	return nil
}
