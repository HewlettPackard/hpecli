// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package log

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetOutput(ioutil.Discard)
	logrus.SetFormatter(&Formatter{NoColors: true})
	logrus.AddHook(&copyHook{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(false)
}

func SetDebugLogging() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
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
		_, _ = h.Stderr.Write([]byte("\n"))
	}

	return nil
}
