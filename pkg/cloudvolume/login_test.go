// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestRunLogins(t *testing.T) {
	const want = "e826d2b3-4925-4f49-86ab-e7f1462c0511"
	jsonResponse := fmt.Sprintf(`{"geo":"US", "token":"%s"}`, want)

	ts := newTestServer("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, jsonResponse)
	})

	opts := cvLoginOptions{
		host:     ts.URL,
		password: "arbitrary",
	}

	// erase value from db - so we know it is empty
	saveData(ts.URL, "")

	err := runLogin(&opts)
	if err != nil {
		t.Fatal(err)
	}

	gotHost, gotToken, err := hostAndToken()
	if err != nil {
		t.Fatal(err)
	}

	if gotHost != ts.URL {
		t.Fatalf(errTempl, gotHost, want)
	}

	if gotToken != want {
		t.Fatalf(errTempl, gotHost, want)
	}
}

func TestRunLoginCmd(t *testing.T) {
	cmd := newLoginCommand()

	// check some fields are set
	if cmd.Use != "login" {
		t.Error("use text not set as expected")
	}

	// just check one of the flags that are set
	if cmd.Flags().Lookup("host") == nil {
		t.Error("didn't find expected flag for host")
	}
}
