// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"reflect"
	"testing"
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
