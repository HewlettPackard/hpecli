// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestLoginWorks(t *testing.T) {
	const want = "e826d2b3-4925-4f49-86ab-e7f1462c0511"
	jsonResponse := fmt.Sprintf(`{"geo":"US", "token":"%s"}`, want)

	ts := newTestServer("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, jsonResponse)
	})

	defer ts.Close()

	// setup login details.  would normally be populated by cobra
	cvLoginData.host = ts.URL

	// erase value from db - so we know it is empty
	storeContext(ts.URL, "")

	err := runCVLogin(nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	d, err := getContext()
	if err != nil {
		t.Fatal(err)
	}

	got := d.APIKey

	if got != want {
		t.Fatalf(errTempl, got, want)
	}
}

func TestHostGetsPrefixed(t *testing.T) {
	ts := newTestServer("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer ts.Close()

	// setup login details.  would normally be populated by cobra
	cvLoginData.host = strings.TrimPrefix(ts.URL, "http://")

	_ = runCVLogin(nil, nil)

	// ensure host got http prefix applied
	if !strings.HasPrefix(cvLoginData.host, "https://") {
		t.Fatal("expected host to get https prefix")
	}
}
