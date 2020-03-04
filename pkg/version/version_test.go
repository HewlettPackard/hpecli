// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package version

import (
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/sirupsen/logrus"
)

const v1 = "0.0.1"
const v0 = "0.0.0"
const expectedError = "got: %v, want: %v"

func TestGetDefault(t *testing.T) {
	want := v0
	got := Get()

	if got != want {
		t.Fatalf(expectedError, got, want)
	}
}

func TestGet(t *testing.T) {
	version = v1
	want := v1
	got := Get()

	if got != want {
		t.Fatalf(expectedError, got, want)
	}
}

func TestGetFull(t *testing.T) {
	version = v1
	gitCommitID = "234a39f"
	buildDate = "2019-01-01"
	want := v1 + ":" + gitCommitID + ":" + buildDate
	got := GetFull()

	if got != want {
		t.Fatalf(expectedError, got, want)
	}
}

func TestCmdCreated(t *testing.T) {
	if Cmd == nil {
		t.Fatal("command should have been initialized")
	}
}

func TestIsFullVersion(t *testing.T) {
	cases := []struct {
		verbose  bool
		want     bool
		logLevel logrus.Level
		name     string
	}{
		{
			name:     "Default short",
			verbose:  false,
			logLevel: logrus.InfoLevel,
			want:     false,
		},
		{
			name:     "verbose is True",
			verbose:  true,
			logLevel: logrus.InfoLevel,
			want:     true,
		},
		{
			name:     "Debug LogLevel is verbose",
			verbose:  false,
			logLevel: logrus.DebugLevel,
			want:     true,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			verbose = c.verbose
			log.Logger.Level = c.logLevel
			got := isFullVersion()
			if got != c.want {
				t.Fatalf(expectedError, got, c.want)
			}
		})
	}
}

func TestFullVersionOutput(t *testing.T) {
	version = v0
	gitCommitID = "0"
	buildDate = "0"

	// if values aren't injected at compile time
	// then everything just defaults to 0
	want := "0.0.0:0:0"
	verbose = true

	got := versionOutput()
	if got != want {
		t.Fatalf(expectedError, got, want)
	}
}

func TestVersionOutput(t *testing.T) {
	// if values aren't injected at compile time
	// then everything just defaults to 0
	want := v0
	verbose = false
	log.Logger.Level = logrus.InfoLevel

	got := versionOutput()
	if got != want {
		t.Fatalf(expectedError, got, want)
	}
}

func TestRun(_ *testing.T) {
	run(nil, nil)
}
