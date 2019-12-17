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
	format          = "%v, %v, %v, all eyes on me!"
	formatExp       = `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.* \[(critical|warning|info|debug)\]\s* \d, \d, \d, all eyes on me!`
	formatWOTimeExp = `\d, \d, \d, all eyes on me!`
)

var (
	a = []interface{}{1, 2, 3}
)

func TestAlways(t *testing.T) {
	TestMode = true
	e, err := regexp.Compile(formatWOTimeExp)
	g := captureLoggerOutput(Always, format, a)

	if err != nil {
		t.Fatalf("Failed to compile regexp '%v': %v", e.String(), err)
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
		t.Fatalf("Failed to compile regexp '%v': %v", e.String(), err)
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
		t.Fatalf("Failed to compile regexp '%v': %v", e.String(), err)
	}

	if !e.MatchString(g) {
		t.Fatalf("Info should produce a pattern '%v' but produces: %v", e.String(), g)
	}
}

func TestInfo(t *testing.T) {
	Level = InfoLevel
	TestMode = true

	e, err := regexp.Compile(formatWOTimeExp)
	g := captureLoggerOutput(Info, format, a)

	if err != nil {
		t.Fatalf("Failed to compile regexp '%v': %v", e.String(), err)
	}

	if !e.MatchString(g) {
		t.Fatalf("Info should produce a pattern '%v' but produces: %v", e.String(), g)
	}
}

func TestDebug(t *testing.T) {
	Level = DebugLevel
	TestMode = true

	e, err := regexp.Compile(formatExp)
	g := captureLoggerOutput(Debug, format, a)

	if err != nil {
		t.Fatalf("Failed to compile regexp '%v': %v", e.String(), err)
	}

	if !e.MatchString(g) {
		t.Fatalf("Info should produce a pattern '%v' but produces: %v", e.String(), g)
	}
}

func captureLoggerOutput(l Logger, format string, a []interface{}) string {
	b := new(bytes.Buffer)
	l(format, append(a, b)...)
	return b.String()
}
