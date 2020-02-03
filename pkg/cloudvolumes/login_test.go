// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolumes


import (
	"testing"
)

func TestCloudVolumeLogin(t *testing.T) {
	
	err := runLogin(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

