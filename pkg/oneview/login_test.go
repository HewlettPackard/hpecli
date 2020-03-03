// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/context"
)

const errTempl = "got: %s, wanted: %s"

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestHostPrefixAdded(t *testing.T) {
	server := newTestServer("/rest/login-sessions", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	ovLoginData.host = strings.Replace(server.URL, "http://", "", 1)
	ovLoginData.password = "blah blah"

	// this will fail with a remote call.. ignore the failure and
	// check the host string to ensure prefix addded
	_ = runOVLogin(nil, nil)

	if !strings.HasPrefix(ovLoginData.host, "https://") {
		t.Fatalf("host should be prefixed with http scheme")
	}
}

func TestAPIKeyIsStored(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer("/rest/login-sessions", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"sessionID":"%s"}`, sessionID)
	})

	defer server.Close()

	ovLoginData.host = server.URL
	ovLoginData.password = "blah blah"

	err := runOVLogin(nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, token, _ := hostAndToken()

	if token != sessionID {
		t.Fatalf(errTempl, token, sessionID)
	}
}

func newTestServer(path string, h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc(path, h)

	return server
}
