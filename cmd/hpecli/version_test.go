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
		verbose  bool
		want     bool
		logLevel logger.LogLevel
		name     string
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
	//if values aren't injected at compile time
	//then everything just defaults to 0
	want := "0.0.0:0:0"
	verbose = true

	got := versionOutput()
	if got != want {
		t.Fatalf(expectedError, got, want)
	}
}

func TestVersionOutput(t *testing.T) {
	//if values aren't injected at compile time
	//then everything just defaults to 0
	want := "0.0.0"
	verbose = false
	logger.Level = logger.InfoLevel

	got := versionOutput()
	if got != want {
		t.Fatalf(expectedError, got, want)
	}
}

func TestRun(_ *testing.T) {
	run(nil, nil)
}
