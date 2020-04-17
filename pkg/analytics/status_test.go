// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewStatusCommand(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{
			name: "get command name for status command",
			want: "status",
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := newStatusCommand(); !reflect.DeepEqual(got.Name(), c.want) {
				t.Errorf("newStatusCommand() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestNewStatusCommandRunE(t *testing.T) {
	cmd := newStatusCommand()

	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckStatus(t *testing.T) {
	cases := []struct {
		cmd            *cobra.Command
		wantErr        bool
		checkAnalytics bool
		name           string
		eventCategory  string
		eventAction    string
	}{
		{
			name:           "check analytics on",
			eventCategory:  "someEC",
			eventAction:    "someEA",
			wantErr:        false,
			checkAnalytics: true,
		},
		{
			name:           "check analytics off",
			eventCategory:  "someEC",
			eventAction:    "someEA",
			wantErr:        false,
			checkAnalytics: false,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if err := checkStatus(c.cmd, c.checkAnalytics, c.eventCategory, c.eventAction); (err != nil) != c.wantErr {
				t.Errorf("checkStatus() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}
