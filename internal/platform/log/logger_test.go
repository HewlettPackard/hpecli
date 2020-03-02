// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package log

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
)

//nolint:funlen // yes.. it is a long test method
func TestHookWritesToCorrectOutput(t *testing.T) {
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
			stdErrWant: `time="0001-01-01T00:00:00Z" level=debug msg="some Message"` + "\n",
		},
		{
			name:       "Info",
			level:      logrus.InfoLevel,
			stdOutWant: "some Message",
			stdErrWant: "",
		},
		{
			name:       "Warn",
			level:      logrus.WarnLevel,
			stdOutWant: "",
			stdErrWant: "some Message",
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			stdOut := bytes.Buffer{}
			stdErr := bytes.Buffer{}

			entry := &logrus.Entry{}
			entry.Logger = logrus.New()
			entry.Logger.Level = test.level
			entry.Message = "some Message"
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

func TestEnsureHookInstalled(t *testing.T) {
	log := New()
	//nolint:gomnd  // don't care about magic number here
	if len(log.Hooks) == 1 {
		t.Errorf("didn't find hook installed")
	}
}
