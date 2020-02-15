// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"errors"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/context"
)

func TestHostPrefixAddedForContext(t *testing.T) {
	iloContextData.host = "127.0.0.1"

	// run it and then check the variable after
	_ = runSetContext(nil, nil)

	if !strings.HasPrefix(iloContextData.host, "https://") {
		t.Fatalf("host should be prefixed with http scheme")
	}
}

func TestContextIsSetInDB(t *testing.T) {
	iloContextData.host = "127.0.0.1"

	// sets the context in the DB
	_ = runSetContext(nil, nil)

	c := iloContext()

	// get value directly from db to verify it was written
	got, _, err := c.APIKey()
	if !errors.Is(err, context.ErrorKeyNotFound) {
		t.Fatal("expected to find the context but not the key")
	}

	if got != iloContextData.host {
		t.Fatal("didn't get expected context after writing")
	}
}
