// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package cfm

import (
	"testing"
)

func TestCFMCommand(t *testing.T) {
	cmd := NewCFMCommand()
	if cmd.Use != "cfm" {
		t.Errorf("command name appears wrong.")
	}
}

func TestCFMSubCommandCount(t *testing.T) {
	const subCommandCount = 4

	cmd := NewCFMCommand()
	if len(cmd.Commands()) != subCommandCount {
		t.Errorf("Wrong number of sub-commands found.")
	}
}
