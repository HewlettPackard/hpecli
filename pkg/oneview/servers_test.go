// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

const shURL = "/rest/server-hardware"

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestAPIKeyPutInServerRequest(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(shURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"sessionID":"%s"}`, sessionID)
		if r.Header.Get("Auth") != sessionID {
			t.Fatal("Expected to find \"auth\" header in request")
		}
	})

	defer server.Close()

	// set context to the test server host
	_ = saveContextAndHostData(server.URL, sessionID)

	// check is above in the http request handler side
	_ = getServerHardware()
}

func TestClientServerRequestFails(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(shURL, func(w http.ResponseWriter, r *http.Request) {
		// cause the request to fail
		w.WriteHeader(http.StatusBadRequest)
	})

	defer server.Close()

	// set context to the test server host
	_ = saveContextAndHostData(server.URL, sessionID)

	// check is above in the http request handler side
	if err := getServers(nil, nil); err == nil {
		t.Fatal("expected to get an error")
	}
}

func TestServerJSONMarshallFails(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(shURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// malformed json will cause the request to fail
		fmt.Fprint(w, `{"type":"server":["bad":"broken"]}`)
	})

	defer server.Close()

	// set context to the test server host
	_ = saveContextAndHostData(server.URL, sessionID)

	// check is above in the http request handler side
	if err := getServerHardware(); err == nil {
		t.Fatal("expected to get an error")
	}
}
