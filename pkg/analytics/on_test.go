// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewOnCommand(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{
			name: "get command name for on command",
			want: "on",
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			print("got \n", newOnCommand())
			if got := newOnCommand(); !reflect.DeepEqual(got.Name(), c.want) {
				t.Errorf("newOnCommand() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestNewOnCommandRunE(t *testing.T) {
	cmd := newOnCommand()

	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}

func TestOnAnalytics(t *testing.T) {
	cases := []struct {
		cmd           *cobra.Command
		enable        bool
		wantErr       bool
		name          string
		eventCategory string
		eventAction   string
	}{
		{
			name:          "check analytics on",
			eventCategory: "someEC",
			eventAction:   "someEA",
			wantErr:       false,
			enable:        true,
		},
		{
			name:          "check analytics off",
			eventCategory: "someEC",
			eventAction:   "someEA",
			wantErr:       false,
			enable:        false,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if err := onAnalytics(c.cmd, c.enable, c.eventCategory, c.eventAction); (err != nil) != c.wantErr {
				t.Errorf("runAnalytics() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}
