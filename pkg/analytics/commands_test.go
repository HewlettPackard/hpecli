// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCommands(t *testing.T) {
	cmd := NewAnalyticsCommand()

	want := 4

	var got int

	// traverse cmd tree to make sure that they all got added
	VisitAll(cmd, func(ccmd *cobra.Command) {
		got++
	})

	if got != want {
		t.Errorf("count of commands doesn't match.  got=%d, want=%d", got, want)
	}
}

func VisitAll(root *cobra.Command, fn func(*cobra.Command)) {
	for _, cmd := range root.Commands() {
		VisitAll(cmd, fn)
	}

	fn(root)
}
