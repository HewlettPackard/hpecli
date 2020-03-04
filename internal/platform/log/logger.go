// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package log

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is the default logger configured to log
// info messages to stdout and everything else to stderr
var Logger = New()

func New() *logrus.Logger {
	return &logrus.Logger{
		Out: ioutil.Discard,
		Formatter: &Formatter{
			NoColors: true,
		},
		Hooks:        addHook(),
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: true,
	}
}

func addHook() logrus.LevelHooks {
	hook := make(logrus.LevelHooks)
	hook.Add(&copyHook{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	return hook
}

type copyHook struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (h *copyHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *copyHook) Fire(entry *logrus.Entry) error {
	// if it info message, write it to stdout
	if entry.Level == logrus.InfoLevel {
		n, _ := h.Stdout.Write([]byte(entry.Message))
		if n > 0 {
			// write \n to match what happens when entry is formatted by entry.String()
			_, _ = h.Stdout.Write([]byte("\n"))
		}

		return nil
	}

	// if we are set for debug logging, write timestamp, level, etc. with msg
	if entry.Logger.GetLevel() == logrus.DebugLevel {
		line, _ := entry.String()
		// this causes formatter.format to be called, which appends \n
		_, _ = h.Stderr.Write([]byte(line))

		return nil
	}

	// just write message
	n, _ := h.Stderr.Write([]byte(entry.Message))
	if n > 0 {
		// write \n to match what happens when entry is formatted by entry.String()
		_, _ = h.Stdout.Write([]byte("\n"))
	}

	return nil
}
