// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"reflect"
	"testing"
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
