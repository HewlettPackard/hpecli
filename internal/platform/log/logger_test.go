// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package log

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
)

//nolint:funlen // yes.. it is a long test method
func TestHookWritesToCorrectOutput(t *testing.T) {
	const msg = "some Message"

	cases := []struct {
		name       string
		level      logrus.Level
		stdOutWant string
		stdErrWant string
	}{
		{
			name:       "Debug",
			level:      logrus.DebugLevel,
			stdOutWant: "",
			stdErrWant: `[0001-01-01T00:00:00Z] [DEBUG] some Message` + "\n",
		},
		{
			name:       "Info",
			level:      logrus.InfoLevel,
			stdOutWant: msg + "\n",
			stdErrWant: "",
		},
		{
			name:       "Warn",
			level:      logrus.WarnLevel,
			stdOutWant: "",
			stdErrWant: msg + "\n",
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			stdOut := bytes.Buffer{}
			stdErr := bytes.Buffer{}

			entry := &logrus.Entry{}
			entry.Logger = logrus.StandardLogger()
			entry.Logger.Level = test.level
			entry.Message = msg
			entry.Level = test.level

			ch := &copyHook{
				Stdout: &stdOut,
				Stderr: &stdErr,
			}

			if err := ch.Fire(entry); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			gotStdOut := stdOut.String()
			gotErrOut := stdErr.String()

			if gotStdOut != test.stdOutWant {
				t.Errorf("error with unexpected data in stdout for level:%v  got=%v -- want=%v",
					test.level, gotStdOut, test.stdOutWant)
			}

			if gotErrOut != test.stdErrWant {
				t.Errorf("error with unexpected data in stderr for level:%v  got=%v -- want=%v",
					test.level, gotErrOut, test.stdErrWant)
			}
		})
	}
}
