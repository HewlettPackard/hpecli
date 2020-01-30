// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package version

import (
	"testing"
)

const v1 = "0.0.1"
const tmplMessage = "Get: got=%s, want=%s"

func TestGetDefault(t *testing.T) {
	want := "0.0.0"
	got := Get()

	if got != want {
		t.Fatalf(tmplMessage, got, want)
	}
}

func TestGet(t *testing.T) {
	version = v1
	want := v1
	got := Get()

	if got != want {
		t.Fatalf(tmplMessage, got, want)
	}
}

func TestGetFull(t *testing.T) {
	version = v1
	gitCommit = "234a39f"
	builtAt = "2019-01-01"
	want := v1 + ":" + gitCommit + ":" + builtAt
	got := GetFull()

	if got != want {
		t.Fatalf(tmplMessage, got, want)
	}
}
