// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package password

import (
	"strings"
	"testing"
)

func TestTwoPasswordOptions(t *testing.T) {
	var p string = "have one"
	err := Read(&p, true, "ilo")

	if err == nil {
		t.Error("expected error")
	}

	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Error("unexpected error text")
	}
}

func TestPasswordAlreadySpecified(t *testing.T) {
	var p string = "have one"
	got := p
	err := Read(&p, false, "ilo")

	if err != nil {
		t.Error("expected error")
	}

	if got != p {
		t.Error("error: password unexpectedly modified")
	}
}
