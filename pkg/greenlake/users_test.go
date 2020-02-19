// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAPIKeyInjectedIntoRequest(t *testing.T) {
	const authValue = "someAuthHeaderValue"

	const tenantID = "someTenantID"

	uriPath := fmt.Sprintf("/scim/v1/tenant/" + tenantID + "/" + "Users")

	server := newTestServer(uriPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth header", r.Header.Get("Authorization"))
		if r.Header.Get("Authorization") != "Bearer "+authValue {
			t.Fatal("Expected to find \"Authorization\" header in request")
		}
		w.WriteHeader(http.StatusOK)
	})

	defer server.Close()

	// set context to the test server host
	storeContext(server.URL, tenantID, authValue)

	_ = runGLGetUsers(nil, nil)
}
