// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCheckCmdCreation(t *testing.T) {
	cases := []struct {
		name     string
		cmd      *cobra.Command
		subCount int
	}{
		{
			name:     "context",
			cmd:      newContextCommand(),
			subCount: 0, //nolint:gomnd // count ok
		},
		{
			name:     "get",
			cmd:      newGetCommand(),
			subCount: 1, //nolint:gomnd // count ok
		},
		{
			name:     "ilo",
			cmd:      NewILOCommand(),
			subCount: 4, //nolint:gomnd // count ok
		},
		{
			name:     "login",
			cmd:      newLoginCommand(),
			subCount: 0, //nolint:gomnd // count ok
		},
		{
			name:     "logout",
			cmd:      newLogoutCommand(),
			subCount: 0, //nolint:gomnd // count ok
		},
		{
			name:     "serviceroot",
			cmd:      newServiceRootCommand(),
			subCount: 0, //nolint:gomnd // count ok
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.cmd.Name() != test.name {
				t.Errorf("command name doesn't match.  got=%v, want=%v",
					test.cmd.Name(), test.name)
			}
			if test.subCount != len(test.cmd.Commands()) {
				t.Errorf("sub command count doesn't match.  got=%v, want=%v",
					len(test.cmd.Commands()), test.subCount)
			}
		})
	}
}
