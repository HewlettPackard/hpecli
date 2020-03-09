// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package log

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
)

const fieldName, fieldValue = "myField", "myValue"

func TestWriteField(t *testing.T) {
	b := bytes.Buffer{}
	entry := &logrus.Entry{}
	entry.Data = make(logrus.Fields, 1)
	entry.Data[fieldName] = fieldValue

	l := Formatter{}
	l.writeField(&b, entry, fieldName)

	got := b.String()
	// Write field formats inside of brackets
	if got != "[myValue] " {
		t.Fatal("didn't retrieve expected value")
	}
}

func TestWriteFields(t *testing.T) {
	b := bytes.Buffer{}
	entry := &logrus.Entry{}
	entry.Data = make(logrus.Fields, 1)
	entry.Data[fieldName] = fieldValue
	entry.Data["otherField"] = "anotherValue"

	l := Formatter{}
	l.writeFields(&b, entry)

	got := b.String()
	// Write field formats inside of brackets
	if got != "[myValue] [anotherValue] " {
		t.Fatal("didn't retrieve expected value")
	}
}

func TestColorLevel(t *testing.T) {
	cases := []struct {
		name  string
		level logrus.Level
		want  int
	}{
		{
			name:  "DebugToGray",
			level: logrus.DebugLevel,
			want:  colorGray,
		},
		{
			name:  "InfoToBlue",
			level: logrus.InfoLevel,
			want:  colorBlue,
		},
		{
			name:  "WarnToYellow",
			level: logrus.WarnLevel,
			want:  colorYellow,
		},
		{
			name:  "ErrorToRed",
			level: logrus.ErrorLevel,
			want:  colorRed,
		},
		{
			name:  "FatalToRed",
			level: logrus.FatalLevel,
			want:  colorRed,
		},
		{
			name:  "PanicToRed",
			level: logrus.PanicLevel,
			want:  colorRed,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := colorForLevel(test.level)
			if got != test.want {
				t.Errorf("didn't get expected value.  got=%v  -  wanted=%v", got, test.want)
			}
		})
	}
}

func TestFormatWorks(t *testing.T) {
	want := "[0001-01-01T00:00:00Z] [PANIC] [myValue] [anotherValue] \n"
	entry := &logrus.Entry{}
	entry.Data = make(logrus.Fields, 1)
	entry.Data[fieldName] = fieldValue
	entry.Data["otherField"] = "anotherValue"

	l := Formatter{NoColors: true}

	b, err := l.Format(entry)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	got := string(b)
	if got != want {
		t.Errorf("unexpected formatted entry.  got=%v -- want=%v", got, want)
	}
}

func TestFormatWithColor(t *testing.T) {
	want := "[0001-01-01T00:00:00Z] \x1b[31m[PANIC] \x1b[0m[myValue] \x1b[0m\n"
	entry := &logrus.Entry{}
	entry.Data = make(logrus.Fields, 1)
	entry.Data[fieldName] = fieldValue

	l := Formatter{}

	b, err := l.Format(entry)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	got := string(b)
	if got != want {
		t.Errorf("unexpected formatted entry.  got=%v -- want=%v", got, want)
	}
}

func TestMarshalLevel(t *testing.T) {
	cases := []struct {
		name  string
		level logrus.Level
	}{
		{
			name:  "TRACE",
			level: logrus.TraceLevel,
		},
		{
			name:  "DEBUG",
			level: logrus.DebugLevel,
		},
		{
			name:  " INFO",
			level: logrus.InfoLevel,
		},
		{
			name:  " WARN",
			level: logrus.WarnLevel,
		},
		{
			name:  "ERROR",
			level: logrus.ErrorLevel,
		},
		{
			name:  "FATAL",
			level: logrus.FatalLevel,
		},
		{
			name:  "PANIC",
			level: logrus.PanicLevel,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := marshalLevel(test.level)
			if got != test.name {
				t.Errorf("didn't get expected value.  got=%v  -  wanted=%v", got, test.name)
			}
		})
	}
}
