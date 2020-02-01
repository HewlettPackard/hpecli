// Copyright © 2017 The Kubicorn Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// This is derived from kris-nova logger: github.com/kris-nova/logger
// It has been  modified to better fit our usage
// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.
package logger

import (
	"bytes"
	"regexp"
	"testing"
)

const (
	format    = "%v, %v, %v, all eyes on me!"
	formatExp = `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.* ` +
		`\[(critical|warning|info|debug)\]\s* \d, \d, \d, all eyes on me!`
	formatWOTimeExp    = `\d, \d, \d, all eyes on me!`
	errorFailedCompile = "Failed to compile regexp '%v': %v"
	errorInfoPattern   = "Info should produce a pattern '%v' but produces: %v"
)

var (
	a = []interface{}{1, 2, 3}
)

func TestAlways(t *testing.T) {
	TestMode = true
	e, err := regexp.Compile(formatWOTimeExp)
	g := captureLoggerOutput(Always, format, a)

	if err != nil {
		t.Fatalf(errorFailedCompile, e.String(), err)
	}

	if !e.MatchString(g) {
		t.Fatalf("Always should produce a pattern '%v' but produces: %v", e.String(), g)
	}
}

func TestCritical(t *testing.T) {
	Level = CriticalLevel
	TestMode = true

	e, err := regexp.Compile(formatWOTimeExp)
	g := captureLoggerOutput(Critical, format, a)

	if err != nil {
		t.Fatalf(errorFailedCompile, e.String(), err)
	}

	if !e.MatchString(g) {
		t.Fatalf("Critical should produce a pattern '%v' but produces: %v", e.String(), g)
	}
}

func TestWarning(t *testing.T) {
	Level = WarningLevel
	TestMode = true

	e, err := regexp.Compile(formatWOTimeExp)
	g := captureLoggerOutput(Warning, format, a)

	if err != nil {
		t.Fatalf(errorFailedCompile, e.String(), err)
	}

	if !e.MatchString(g) {
		t.Fatalf(errorInfoPattern, e.String(), g)
	}
}

func TestInfo(t *testing.T) {
	Level = InfoLevel
	TestMode = true

	e, err := regexp.Compile(formatWOTimeExp)
	g := captureLoggerOutput(Info, format, a)

	if err != nil {
		t.Fatalf(errorFailedCompile, e.String(), err)
	}

	if !e.MatchString(g) {
		t.Fatalf(errorInfoPattern, e.String(), g)
	}
}

func TestDebug(t *testing.T) {
	Level = DebugLevel
	TestMode = true

	e, err := regexp.Compile(formatExp)
	g := captureLoggerOutput(Debug, format, a)

	if err != nil {
		t.Fatalf(errorFailedCompile, e.String(), err)
	}

	if !e.MatchString(g) {
		t.Fatalf(errorInfoPattern, e.String(), g)
	}
}

func TestLogLevelToStrings(t *testing.T) {
	cases := []struct {
		name string
		t    LogLevel
	}{
		{
			name: "debug",
			t:    DebugLevel,
		},
		{
			name: "info",
			t:    InfoLevel,
		},
		{
			name: "warning",
			t:    WarningLevel,
		},
		{
			name: "critical",
			t:    CriticalLevel,
		},
		{
			name: "always",
			t:    AlwaysLevel,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got := c.t.String()
			want := c.name
			if got != want {
				t.Fatalf("Didn't get the expected result.  got=%s, want=%s", got, want)
			}
		})
	}
}

func TestSetLogLevel(t *testing.T) {
	cases := []struct {
		name string
		t    LogLevel
	}{
		{
			name: "debug",
			t:    DebugLevel,
		},
		{
			name: "info",
			t:    InfoLevel,
		},
		{
			name: "warning",
			t:    WarningLevel,
		},
		{
			name: "critical",
			t:    CriticalLevel,
		},
		{
			name: "always",
			t:    AlwaysLevel,
		},
		{
			// defaults to warning
			name: "unknown",
			t:    WarningLevel,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			SetLogLevel(c.name)
			got := Level
			want := c.t
			if got != want {
				t.Fatalf("Didn't set the expected result.  got=%s, want=%s", got, want)
			}
		})
	}
}

func TestColorWritesStdOut(t *testing.T) {
	Level = CriticalLevel
	TestMode = false

	e, err := regexp.Compile(formatWOTimeExp)
	g := captureLoggerOutput(Critical, format, a)

	if err != nil {
		t.Fatalf(errorFailedCompile, e.String(), err)
	}

	if g != "" {
		t.Fatal("output should not be capture.  output should be written to stdout")
	}
}

func captureLoggerOutput(l Logger, format string, a []interface{}) string {
	b := new(bytes.Buffer)
	l(format, append(a, b)...)

	return b.String()
}
