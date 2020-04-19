// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"reflect"
	"testing"
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
