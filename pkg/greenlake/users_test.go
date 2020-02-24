// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestGLAPIKeyInjectedIntoRequest(t *testing.T) {
	const authValue = "someAuthorizationHeaderValue"

	const tenantID = "someTenantID"

	uriPath := fmt.Sprintf("/scim/v1/tenant/" + tenantID + "/" + "Users")

	server := newTestServer(uriPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+authValue {
			t.Fatal("Expected to find \"Authorization\" header in request")
		}
		w.WriteHeader(http.StatusOK)
	})

	defer server.Close()

	// set context to the test server host
	saveData(server.URL, tenantID, authValue)

	_ = runGLGetUsers(nil, nil)
}
