// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"testing"
)

func TestCheckCmdCreation(t *testing.T) {
	cmd := NewCloudVolumeCommand()

	if cmd.Name() != "cloudvolume" {
		t.Error("name not set on command")
	}

	if len(cmd.Commands()) != 2 { //nolint:gomnd  // number ok here
		t.Error("unexpected discrepency in sub command count")
	}
}
