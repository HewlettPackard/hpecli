// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package log

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
)

// Formatter - logrus formatter, implements logrus.Formatter
type Formatter struct {
	NoColors bool // disable colors
}

// Format an log entry
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := colorForLevel(entry.Level)

	// output buffer
	b := &bytes.Buffer{}

	// write time
	fmt.Fprintf(b, "[%s] ", entry.Time.Format(time.RFC3339))

	// write level
	level := marshalLevel(entry.Level)

	if !f.NoColors {
		fmt.Fprintf(b, "\x1b[%dm", levelColor)
	}

	fmt.Fprintf(b, "[%s] ", level)

	if !f.NoColors {
		b.WriteString("\x1b[0m")
	}

	// write fields
	f.writeFields(b, entry)

	if entry.HasCaller() {
		fmt.Fprintf(b, "[%s():%d] - ", path.Base(entry.Caller.Function), entry.Caller.Line)
	}

	if !f.NoColors {
		b.WriteString("\x1b[0m")
	}

	// write message
	b.WriteString(entry.Message)

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *Formatter) writeFields(b io.Writer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeField(b io.Writer, entry *logrus.Entry, field string) {
	fmt.Fprintf(b, "[%v] ", entry.Data[field])
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func colorForLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}

func marshalLevel(level logrus.Level) string {
	switch level {
	case logrus.TraceLevel:
		return "TRACE"
	case logrus.DebugLevel:
		return "DEBUG"
	case logrus.InfoLevel:
		return " INFO"
	case logrus.WarnLevel:
		return " WARN"
	case logrus.ErrorLevel:
		return "ERROR"
	case logrus.FatalLevel:
		return "FATAL"
	case logrus.PanicLevel:
		return "PANIC"
	}

	return "UNKNOWN"
}
