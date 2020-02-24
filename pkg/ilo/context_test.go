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
	const h1 = "host1"

	const v1 = "value1"

	if err := saveData(h1, v1); err != nil {
		t.Fatal(err)
	}

	gotHost, gotToken, err := hostAndToken()
	if err != nil {
		t.Fatal(err)
	}

	if gotHost != h1 {
		t.Fatal("didn't retrieve matching host")
	}

	if gotToken != v1 {
		t.Fatal("didn't retrieve matching value")
	}
}
