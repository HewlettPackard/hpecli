// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestGLHostPrefixAddedForContext(t *testing.T) {
	// clear everything from the mock store
	context.MockClear()

	host := "127.0.0.2"
	want := "https://" + host

	cmd := newContextCommand()
	cmd.SetArgs([]string{"--host", host})

	_ = cmd.Execute()

	// see if the naked host got stored (it shouldn't)
	got, _ := getContext()
	if got != want {
		t.Error("didn't find context with corrected host: " + want)
	}
}

func TestGLGetWorks(t *testing.T) {
	const h1, t1, tn1 = "host1", "token1", "tenant1"

	d := &sessionData{h1, t1, tn1}

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

	if got.TenantID != tn1 {
		t.Fatal("didn't retrieve matching tenant")
	}
}
