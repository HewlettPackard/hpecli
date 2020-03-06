// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"errors"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestHostPrefixAddedForContext(t *testing.T) {
	host := "127.0.1.1"

	cmd := newContextCommand()
	cmd.SetArgs([]string{"--host", host})
	_ = cmd.Execute()

	// check the db to make sure it was persisted
	got, err := getContext()
	if err != nil {
		t.Fatal(err)
	}

	if got != "https://"+host {
		t.Fatal("context didn't get set correctly")
	}
}

func TestContextIsSetInDB(t *testing.T) {
	host := "https://127.0.0.1"

	// sets the context in the DB
	_ = runSetContext(host)

	_, _, err := hostAndToken()
	// we have set the context but not the apikey value
	// so we should get a ErrorKeyNotFound
	if !errors.Is(err, context.ErrorKeyNotFound) {
		t.Fatal("expected to find the context but not the key")
	}
}

func TestCheckDefaultContextFound(t *testing.T) {
	// setup data
	setContext("https://127.0.0.1")

	host := ""
	if err := runSetContext(host); err != nil {
		t.Errorf("didn't get default context successfully: %s", err)
	}
}
