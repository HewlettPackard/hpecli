// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestHostPrefixAddedForContext(t *testing.T) {
	iloContextHost.host = "127.0.0.1"

	// run it and then check the variable after
	_ = runSetContext(nil, nil)

	if !strings.HasPrefix(iloContextHost.host, "https://") {
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
