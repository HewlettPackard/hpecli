// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cfm

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

const logoutURI = "/api/v1/auth/token"
const expectedErrMsg = "expected to see an error here but didn't"

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestLogoutRequestPasses(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(logoutURI, func(w http.ResponseWriter, r *http.Request) {
		// cause the request to pass
		w.WriteHeader(http.StatusOK)
	})

	defer server.Close()

	// set context to the test server host
	_ = saveContextAndHostData(strings.Replace(server.URL, "https://", "", -1), sessionID)

	// check is above in the http request handler side
	if err := runLogout("host"); err == nil {
		t.Fatal(expectedErrMsg)
	}
}

func TestLogoutRequestFails(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(logoutURI, func(w http.ResponseWriter, r *http.Request) {
		// cause the request to fail
		w.WriteHeader(http.StatusBadRequest)
	})

	defer server.Close()

	// set context to the test server host
	_ = saveContextAndHostData(strings.Replace(server.URL, "https://", "", -1), sessionID)

	// check is above in the http request handler side
	if err := runLogout("host"); err == nil {
		t.Fatal(expectedErrMsg)
	}
}

func TestLogoutRemovesAPIKeyFromContext(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(logoutURI, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	defer server.Close()

	// set context to the test server host
	_ = saveContextAndHostData(strings.Replace(server.URL, "https://", "", -1), sessionID)

	_ = runLogout("host")

	// verify the data is gone
	var token string

	c := context.New(CFMContextKey)

	err := c.HostData(dataKey(server.URL), &token)
	if !errors.Is(err, context.ErrorKeyNotFound) {
		t.Fatal("expected ErrorKeyNotFound error")
	}
}

func TestLogoutRemovesAPIKeyFromParameter(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(logoutURI, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	defer server.Close()

	_ = saveContextAndHostData(strings.Replace(server.URL, "https://", "", -1), sessionID)
	// erase the context .. just to make sure that it doesn't pickup the host from the context
	_ = setContext("")

	// specify the param like it was passed on the command line
	_ = runLogout(server.URL)

	// verify the data is gone
	var token string

	c := context.New(CFMContextKey)

	err := c.HostData(dataKey(server.URL), &token)
	if !errors.Is(err, context.ErrorKeyNotFound) {
		t.Fatal("expected ErrorKeyNotFound error")
	}
}

func TestLogoutFailsWhenItCantGetContext(t *testing.T) {
	// erase context for runnig the command
	_ = setContext("")

	// run the command that will fail because of the missing context
	if err := runLogout("host"); err == nil {
		t.Fatal(expectedErrMsg)
	}
}
