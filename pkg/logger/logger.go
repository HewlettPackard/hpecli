// Package logger is used by the CLI to write information to the
// console.  Debug level logging will show timestamps and log level
// in addition to specific debug statements.  It will also write
// critical messages in Red if the console supports it
//
// Copyright © 2017
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
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Logger interface .. enables testing
type Logger func(format string, a ...interface{})

// LogLevel needs doc
type LogLevel int

const (
	// AlwaysLevel need doc
	AlwaysLevel LogLevel = iota
	// CriticalLevel need doc
	CriticalLevel
	// WarningLevel need doc
	WarningLevel
	// InfoLevel need doc
	InfoLevel
	// DebugLevel need doc
	DebugLevel
)

// Convert the Level to a string. E.g. AlwaysLevel becomes "always".
func (level LogLevel) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarningLevel:
		return "warning"
	case CriticalLevel:
		return "critical"
	case AlwaysLevel:
		return "always"
	}
	return "unknown"
}

var (
	// Level Won't log below this level
	Level = WarningLevel
	// Color determines if we use color console
	Color = true
	// TestMode disables color if in test mode
	TestMode = false
)

// Always writes output to the console
func Always(format string, a ...interface{}) {
	write(AlwaysLevel, format, a...)
}

// Critical writes output to the console
func Critical(format string, a ...interface{}) {
	write(CriticalLevel, format, a...)
}

// Warning writes output to the console
func Warning(format string, a ...interface{}) {
	write(WarningLevel, format, a...)
}

// Info writes output to the console
func Info(format string, a ...interface{}) {
	write(InfoLevel, format, a...)
}

// Debug writes output to the console
func Debug(format string, a ...interface{}) {
	write(DebugLevel, format, a...)
}

// SetLogLevel based on text string
func SetLogLevel(logLevel string) {
	switch strings.ToLower(logLevel) {
	case AlwaysLevel.String():
		Level = AlwaysLevel
	case CriticalLevel.String():
		Level = CriticalLevel
	case WarningLevel.String(), "warn":
		Level = WarningLevel
	case InfoLevel.String():
		Level = InfoLevel
	case DebugLevel.String():
		Level = DebugLevel
	default:
		Level = WarningLevel
		Warning("%s is not a valid logging level.  Defaulting to \"warn\"", logLevel)
	}
}

func write(lvl LogLevel, format string, a ...interface{}) {
	if Level >= lvl {
		a, w := extractLoggerArgs(format, a...)
		if !strings.Contains(format, "\n") {
			format = fmt.Sprintf("%s%s", format, "\n")
		}
		if Level >= DebugLevel {
			format = prefixDebug(format, lvl)
		}
		s := fmt.Sprintf(format, a...)

		if Color {
			if CriticalLevel == lvl {
				s = color.RedString(s)
			} else {
				s = color.WhiteString(s)
			}
			if !TestMode {
				w = color.Output
			}

		}
		fmt.Fprintf(w, s)
	}
}

func extractLoggerArgs(format string, a ...interface{}) ([]interface{}, io.Writer) {
	var w io.Writer = os.Stdout
	if n := len(a); n > 0 {
		// extract an io.Writer at the end of a
		if value, ok := a[n-1].(io.Writer); ok {
			w = value
			a = a[0 : n-1]
		}
	}

	return a, w
}

func prefixDebug(format string, level LogLevel) string {
	t := time.Now()
	rfct := t.Format(time.RFC3339)
	return fmt.Sprintf("%s [%s]  %s", rfct, level, format)
}
