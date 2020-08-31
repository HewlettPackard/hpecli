// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package root

import (
	"testing"
	"github.com/HewlettPackard/hpecli/pkg/version"
)



func TestNewRoot(t *testing.T) {

	vInfo := &version.Info{
		Sematic:     "",
		GitCommitID: "",
		BuildDate:   "",
		Verbose:     false,
	}

	root := NewRootCommand(vInfo)

	
	root.Commands()
}
