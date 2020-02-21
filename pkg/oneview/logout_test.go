// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"net/http"
	"testing"
)

const logoutURI = "/rest/login-sessions"
const expectedErrMsg = "expected to see an error here but didn't"

func TestLogoutRequestFails(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(logoutURI, func(w http.ResponseWriter, r *http.Request) {
		// cause the request to fail
		w.WriteHeader(http.StatusBadRequest)
	})

	defer server.Close()

	// set context to the test server host
	_ = storeContext(server.URL, sessionID)

	// check is above in the http request handler side
	if err := runOVLogout(nil, nil); err == nil {
		t.Fatal(expectedErrMsg)
	}
}

func TestLogoutRemovesContext(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(logoutURI, func(w http.ResponseWriter, r *http.Request) {
		// cause the request to fail
		w.WriteHeader(http.StatusOK)
	})

	defer server.Close()

	// set context to the test server host
	_ = storeContext(server.URL, sessionID)

	// check is above in the http request handler side
	_ = runOVLogout(nil, nil)

	if _, err := getContext(); err == nil {
		t.Fatal(expectedErrMsg)
	}
}

func TestLogoutFailsWhenItCantGetContext(t *testing.T) {
	// erase context for runnig the command
	_ = removeContext()

	// run the command that will fail because of the missing context
	if err := runOVLogout(nil, nil); err == nil {
		t.Fatal(expectedErrMsg)
	}
}
