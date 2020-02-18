// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"net/http"
	"testing"
)

const rootURL = "/redfish/v1/"

func TestAPIKeyInjectedIntoRequest(t *testing.T) {
	const authValue = "someAuthHeaderValue"

	server := newTestServer(rootURL, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Auth-Token") != authValue {
			t.Fatal("Expected to find \"auth\" header in request")
		}
		w.WriteHeader(http.StatusOK)
	})

	defer server.Close()

	// set context to the test server host
	storeContext(server.URL, authValue)

	_ = runILOServiceRoot(nil, nil)
}
