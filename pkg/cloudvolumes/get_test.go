// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolumes


import (
	"testing"
)

func TestCloudVolumeGet(t *testing.T) {
	
	err := runGet(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}
