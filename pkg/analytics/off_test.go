// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewOffCommand(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{
			name: "get command name for off command",
			want: "off",
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := newOffCommand(); !reflect.DeepEqual(got.Name(), c.want) {
				t.Errorf("newOffCommand() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestNewOffCommandRunE(t *testing.T) {
	cmd := newOffCommand()

	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}

func TestOffAnalytics(t *testing.T) {
	cases := []struct {
		cmd           *cobra.Command
		disable       bool
		wantErr       bool
		name          string
		eventCategory string
		eventAction   string
	}{
		{
			name:          "check analytics off",
			eventCategory: "someEC",
			eventAction:   "someEA",
			wantErr:       false,
			disable:       true,
		},
		{
			name:          "check analytics on",
			eventCategory: "someEC",
			eventAction:   "someEA",
			wantErr:       false,
			disable:       false,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if err := offAnalytics(c.cmd, c.disable, c.eventCategory, c.eventAction); (err != nil) != c.wantErr {
				t.Errorf("offAnalytics() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}
