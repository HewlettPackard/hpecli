// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package cfm

import (
	"testing"
)

func TestCFMGetCommand(t *testing.T) {
	cmd := newGetCommand()
	if cmd.Use != "get" {
		t.Errorf("command name appears wrong.")
	}
}

func TestCFMGetSubCommandCount(t *testing.T) {
	const subCommandCount = 4

	cmd := NewCFMCommand()
	if len(cmd.Commands()) != subCommandCount {
		t.Errorf("Wrong number of sub-commands found.")
	}
}
