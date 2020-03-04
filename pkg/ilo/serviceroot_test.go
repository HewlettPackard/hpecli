// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"net/http"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

const rootURL = "/redfish/v1/"

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestAPIKeyInjectedIntoRequest(t *testing.T) {
	const authValue = "someAuthHeaderValue"

	server := newTestServer(rootURL, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Auth-Token") != authValue {
			t.Fatal("Expected to find \"auth\" header in request")
		}
		w.WriteHeader(http.StatusOK)
	})

	defer server.Close()

	d := &sessionData{server.URL, authValue, ""}

	// set context to the test server host
	saveContextAndSessionData(d)

	_ = runILOServiceRoot(nil, nil)
}
