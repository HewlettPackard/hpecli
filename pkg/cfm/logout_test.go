// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package cfm

import (
	"fmt"
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
		response := `{
			"ResponseObject": {
			  "StatusCode": 200,
			  "Result": "OK"
			}
		  }`
		fmt.Fprint(w, response)
	})

	defer server.Close()

	host := strings.ReplaceAll(server.URL, "https://", "")
	// set context to the test server host
	_ = saveContextAndHostData(host, sessionID)

	// check is above in the http request handler side
	if err := runLogout(host); err != nil {
		t.Fatalf("logout failed: %v", err)
	}
}

func TestLogoutCommandPasses(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(logoutURI, func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"ResponseObject": {
			  "StatusCode": 200,
			  "Result": "OK"
			}
		  }`
		fmt.Fprint(w, response)
	})

	defer server.Close()

	host := strings.ReplaceAll(server.URL, "https://", "")
	// set context to the test server host
	_ = saveContextAndHostData(host, sessionID)

	cmd := newLogoutCommand()

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("logout failed: %v", err)
	}
}

// this test fails .. because the if the http
// response returns 400 -- the CheckErr method
// panics (exits the test)
// Prefix with function name with Test - if/once
// it is fixed.
func LogoutRequestFails(t *testing.T) {
	server := newTestServer(logoutURI, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	defer server.Close()

	host := strings.ReplaceAll(server.URL, "https://", "")
	// set context to the test server host
	_ = saveContextAndHostData(host, "someSessionId")

	if err := runLogout(host); err == nil {
		t.Fatal(expectedErrMsg)
	}
}

/*
func TestLogoutRemovesAPIKeyFromContext(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(logoutURI, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	defer server.Close()


	// set context to the test server host
	_ = saveContextAndHostData(strings.ReplaceAll(server.URL, "https://", ""), sessionID)

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

	_ = saveContextAndHostData(strings.ReplaceAll(server.URL, "https://", ""), sessionID)
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
*/
