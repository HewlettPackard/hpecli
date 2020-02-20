// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"errors"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestHostPrefixAddedForContext(t *testing.T) {
	ovContextHost.host = "127.0.0.1"

	// run it and then check the variable after
	_ = runChangeContext(nil, nil)

	if !strings.HasPrefix(ovContextHost.host, "https://") {
		t.Fatalf("host should be prefixed with http scheme")
	}
}

func TestContextIsSetInDB(t *testing.T) {
	ovContextHost.host = "https://127.0.0.1"

	// sets the context in the DB
	_ = runChangeContext(nil, nil)

	_, err := getContext()
	// we have set the context but not the apikey value
	// so we should get a ErrorKeyNotFound
	if !errors.Is(err, context.ErrorKeyNotFound) {
		t.Fatal("expected to find the context but not the key")
	}
}
