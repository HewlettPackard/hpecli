// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newOffCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "off",
		Short: "Turn off Google Analytics",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return offAnalytics(cmd)
		},
	}

	return cmd
}

func offAnalytics(cmd *cobra.Command) error {
	fmt.Printf("inside analytics OFF before %s \n", os.Getenv("HPECLI_DISABLE_ANALYTICS"))
	os.Setenv("HPECLI_DISABLE_ANALYTICS", "true")
	fmt.Printf("inside analytics OFF after %s \n", os.Getenv("HPECLI_DISABLE_ANALYTICS"))
	analytics := NewAnalyticsClient("1", "event", "analytics", "off", "200", "", "hpecli/0.0.1", "0.0.1", "hpecli")
	// return errors.New("can't work with analytics")

	resp, err := analytics.TrackEvent()
	if err != nil {
		logrus.Warningf("Unable to send the analytics info with supplied event details")
		return err
	}
	fmt.Printf("Successssssss analytics ON %s \n ", resp)
	return nil
}
