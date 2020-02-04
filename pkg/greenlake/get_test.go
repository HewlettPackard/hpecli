//(C) Copyright 2019 Hewlett Packard Enterprise Development LP

package greenlake

import (
	"fmt"
	"net/http"
	"testing"
)

const userURL = "/scim/v1/tenant/dummy_tenent_id/Users"

func TestClientRequestFails(t *testing.T) {
	const accessToken = "HERE_IS_A_ID"
	const tenantID = "dummy_tenent_id"
	getPath = "users"
	getJSONResult = false

	server := newTestServer(userURL, func(w http.ResponseWriter, r *http.Request) {
		// cause the request to fail
		w.WriteHeader(http.StatusBadRequest)
	})

	defer server.Close()

	// set context to the test server host
	_ = setTokenTenantID(server.URL, tenantID, accessToken)

	// check is above in the http request handler side
	if err := runGlGet(nil, nil); err == nil {
		t.Fatal("expected to get an error")
	}
}

func TestPrintJSONResponse(t *testing.T) {
	const accessToken = "HERE_IS_A_ID"
	const tenantID = "dummy_tenent_id"
	getPath = "users"
	getJSONResult = true

	server := newTestServer(userURL, func(w http.ResponseWriter, r *http.Request) {
		// cause the request to fail
		w.WriteHeader(http.StatusBadRequest)
	})

	defer server.Close()

	// set context to the test server host
	_ = setTokenTenantID(server.URL, tenantID, accessToken)

	// check is above in the http request handler side
	if err := runGlGet(nil, nil); err == nil {
		t.Fatal("expected to get an error")
	}
}

func TestJSONMarshallFails(t *testing.T) {
	const accessToken = "HERE_IS_A_ID"
	const tenantID = "dummy_tenent_id"
	getPath = "users"
	getJSONResult = false

	server := newTestServer(userURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// malformed json will cause the request to fail
		fmt.Fprint(w, `{"type":"server":["bad":"broken"]}`)
	})

	defer server.Close()

	// set context to the test server host
	_ = setTokenTenantID(server.URL, tenantID, accessToken)

	// check is above in the http request handler side
	if err := runGlGet(nil, nil); err == nil {
		t.Fatal("expected to get an error")
	}
}
