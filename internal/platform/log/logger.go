// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package log

import (
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
	hook.Add(&CopyHook{})

	return hook
}

type CopyHook struct {
}

func (h *CopyHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *CopyHook) Fire(entry *logrus.Entry) error {
	// if it info message, write it to stdout
	if entry.Level == logrus.InfoLevel {
		os.Stdout.Write([]byte(entry.Message))
		return nil
	}

	// if we are set for debug logging, write timestamp, level, etc. with msg
	if entry.Logger.GetLevel() == logrus.DebugLevel {
		line, _ := entry.String()
		os.Stderr.Write([]byte(line))

		return nil
	}

	// just write message
	os.Stderr.Write([]byte(entry.Message))

	return nil
}
