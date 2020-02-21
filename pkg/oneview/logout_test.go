// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"net/http"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestClientLogoutRequestFails(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(shURL, func(w http.ResponseWriter, r *http.Request) {
		// cause the request to fail
		w.WriteHeader(http.StatusBadRequest)
	})

	defer server.Close()

	// set context to the test server host
	_ = storeContext(server.URL, sessionID)

	// check is above in the http request handler side
	if err := runOVLogout(nil, nil); err == nil {
		t.Fatal("expected to get an error")
	}
}
