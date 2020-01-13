// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package main

import (
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/logger"
)

const expectedError = "got: %v, want: %v"

func TestCmdCreated(t *testing.T) {
	if versionCmd == nil {
		t.Fatal("command should have been initialized")
	}
}

func TestIsFullVersion(t *testing.T) {

	cases := []struct {
		name     string
		verbose  bool
		logLevel logger.LogLevel
		want     bool
	}{
		{
			name:     "Default short",
			verbose:  false,
			logLevel: logger.InfoLevel,
			want:     false,
		},
		{
			name:     "verbose is True",
			verbose:  true,
			logLevel: logger.InfoLevel,
			want:     true,
		},
		{
			name:     "Debug LogLevel is verbose",
			verbose:  false,
			logLevel: logger.DebugLevel,
			want:     true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			verbose = c.verbose
			logger.Level = c.logLevel
			got := isFullVersion()
			if got != c.want {
				t.Fatalf(expectedError, got, c.want)
			}
		})
	}

}

func TestFullVersionOutput(t *testing.T) {
	//This is hacky... because values aren't
	//injected, then we just end up with
	//colon seperators
	want := "::"
	verbose = true

	got := versionOutput()
	if got != want {
		t.Fatalf(expectedError, got, want)
	}
}

func TestVersionOutput(t *testing.T) {
	//This is hacky... because values aren't
	//injected, then we just end up with
	//empty return
	want := ""
	verbose = false
	logger.Level = logger.InfoLevel

	got := versionOutput()
	if got != want {
		t.Fatalf(expectedError, got, want)
	}
}
