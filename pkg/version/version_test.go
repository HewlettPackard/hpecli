// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package version

import (
	"testing"
)

const incorrectValueError = "incorrect value: got=%v, want=%v"

func TestVersionInfoString(t *testing.T) {
	cases := []struct {
		name string
		want string
		vi   *Info
	}{
		{
			name: "NotVerbose",
			vi: &Info{
				BuildDate:   "A",
				GitCommitID: "B",
				Sematic:     "C",
				Verbose:     false,
			},
			want: "C",
		},
		{
			name: "Verbose",
			vi: &Info{
				BuildDate:   "A",
				GitCommitID: "B",
				Sematic:     "C",
				Verbose:     true,
			},
			want: "C:B:A",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got := c.vi.String()
			if got != c.want {
				t.Errorf(incorrectValueError, got, c.want)
			}
		})
	}
}

func TestVerboseFromCommandLine(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want bool
		vi   *Info
	}{
		{
			name: "NotVerbose",
			args: []string{"-v"},
			vi: &Info{
				Verbose: false,
			},
			want: true,
		},
		{
			name: "Verbose",
			args: []string{"--verbose"},
			vi: &Info{
				Verbose: false,
			},
			want: true,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			cmd := NewVersionCommand(c.vi)

			cmd.SetArgs(c.args)
			cmd.Execute()

			got := c.vi.Verbose

			if got != c.want {
				t.Errorf(incorrectValueError, got, c.want)
			}
		})
	}
}
