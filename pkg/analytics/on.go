// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// EnvDisableAnalytics is environmental variable to disable google analytics
const EnvDisableAnalytics = "HPECLI_DISABLE_ANALYTICS"

func newOnCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "on",
		Short: "Turn on Google Analytics",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runAnalytics(cmd)
		},
	}

	return cmd
}

func runAnalytics(cmd *cobra.Command) error {
	fmt.Printf("inside analytics ON %s \n ", os.Getenv("HPECLI_DISABLE_ANALYTICS"))
	os.Unsetenv("HPECLI_DISABLE_ANALYTICS")
	analytics := NewAnalyticsClient("1", "event", "analytics", "on", "200", "", "hpecli/0.0.1", "0.0.1", "hpecli")
	// return errors.New("can't work with analytics")
	resp, err := analytics.TrackEvent()
	if err != nil {
		logrus.Warningf("Unable to send the analytics info with supplied event details")
		return err
	}
	fmt.Printf("Successssssss analytics ON %s \n ", resp)
	return nil
}
