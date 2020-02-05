// (C) Copyright 2019 Hewlett Packard Enterprise Development LP

package greenlake

import (
	"fmt"
	"net/http"
	"testing"
)

const userURL = "/scim/v1/tenant/dummy_tenent_id/Users"
const errExp = "expected to get an error"

func TestGLClientRequestFails(t *testing.T) {
	const accessToken = "GreenLake_Access_Token"

	const tenantID = "dummy_tenent_id"

	getPath = users

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
		t.Fatal(errExp)
	}
}

func TestPrintJSONResponse(t *testing.T) {
	const accessToken = "GreenLake_Access_Token"

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
		t.Fatal(errExp)
	}
}

func TestGLJSONMarshallFails(t *testing.T) {
	const accessToken = "GreenLake_Access_Token"

	const tenantID = "dummy_tenent_id"

	getPath = "users"

	getJSONResult = false

	server := newTestServer(userURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// malformed json will cause the request to fail
		fmt.Fprint(w, `{"type":"glserver":["bad":"broken"]}`)
	})

	defer server.Close()

	// set context to the test server host
	_ = setTokenTenantID(server.URL, tenantID, accessToken)

	// check is above in the http request handler side
	if err := runGlGet(nil, nil); err == nil {
		t.Fatal(errExp)
	}
}
