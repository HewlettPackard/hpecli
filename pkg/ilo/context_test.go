// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestHostPrefixAddedForContext(t *testing.T) {
	host := "127.0.0.1"

	// run it and then check the variable after
	_ = runSetContext(&host)

	if !strings.HasPrefix(host, "https://") {
		t.Fatalf("host should be prefixed with http scheme")
	}
}

func TestGetWorks(t *testing.T) {
	const h1, t1, l1 = "host1", "token1", "location1"

	d := &sessionData{h1, t1, l1}

	if err := saveContextAndSessionData(d); err != nil {
		t.Fatal(err)
	}

	got, err := defaultSessionData()
	if err != nil {
		t.Fatal(err)
	}

	if got.Host != h1 {
		t.Fatal("didn't retrieve matching host")
	}

	if got.Token != t1 {
		t.Fatal("didn't retrieve matching token")
	}

	if got.Location != l1 {
		t.Fatal("didn't retrieve matching value")
	}
}

func TestNewContextCommand(t *testing.T) {
	cmd := newContextCommand()

	if cmd.Use != "context" {
		t.Error("unexpected use value")
	}
}

func TestCheckDefaultContextFound(t *testing.T) {
	// setup data
	setContext("https://127.0.0.1")

	host := ""
	// don't specify a host, so it will look for the default context value
	if err := runSetContext(&host); err != nil {
		t.Errorf("didn't get default context successfully: %s", err)
	}
}
