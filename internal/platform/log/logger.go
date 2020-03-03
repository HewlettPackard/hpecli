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
		_, _ = h.Stdout.Write([]byte(entry.Message))

		return nil
	}

	// if we are set for debug logging, write timestamp, level, etc. with msg
	if entry.Logger.GetLevel() == logrus.DebugLevel {
		line, _ := entry.String()
		_, _ = h.Stderr.Write([]byte(line))

		return nil
	}

	// just write message
	_, _ = h.Stderr.Write([]byte(entry.Message))

	return nil
}
