// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHostPrefixAdded(t *testing.T) {
	server := newTestServer("/identity/v1/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	glLoginData.host = strings.Replace(server.URL, "http://", "", 1)

	// this will fail with a remote call.. ignore the failure and
	// check the host string to ensure prefix addded
	_ = runGLLogin(nil, nil)

	if !strings.HasPrefix(glLoginData.host, "http://") {
		t.Fatalf("host should be prefixed with http scheme")
	}
}

func TestAPIKeyIsStored(t *testing.T) {
	const accessToken = "HERE_IS_A_ID"
	// const tenantId = "HERE_IS_A_TenantID"

	server := newTestServer("/identity/v1/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"access_token":"%s"}`, accessToken)
	})

	defer server.Close()

	glLoginData.host = server.URL

	err := runGLLogin(nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// sessionId is stored - so get it and verify it
	_, _, got := getTokenTenantID()
	if got != accessToken {
		t.Fatal(errTempl, got, accessToken)
	}
}

func newTestServer(path string, h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc(path, h)

	return server
}
