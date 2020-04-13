// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cfm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

const errTempl = "got: %s, wanted: %s"

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestLogin(t *testing.T) {
	server := newTestServer("/api/v1/auth/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"count":1, "result":"some_uuid_token"}`)
	})

	defer server.Close()

	opts := &cfmLoginOptions{
		host:     strings.ReplaceAll(server.URL, "https://", ""),
		username: "admin",
		password: "somePassword",
	}

	err := runLogin(opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoginFails(t *testing.T) {
	server := newTestServer("/api/v1/auth/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"count":1, "result":"some error message"}`)
	})

	defer server.Close()

	opts := &cfmLoginOptions{
		host:     strings.ReplaceAll(server.URL, "https://", ""),
		username: "admin",
		password: "somePassword",
	}

	err := runLogin(opts)
	if err == nil {
		t.Fatal(err)
	}
}

func TestAPIKeyIsStored(t *testing.T) {
	const sessionID = "SomeSessionID"

	server := newTestServer("/api/v1/auth/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"sessionID":"%s"}`, sessionID)
	})

	defer server.Close()

	opts := &cfmLoginOptions{
		host:     strings.ReplaceAll(server.URL, "https://", ""),
		password: "blah blah",
	}

	err := runLogin(opts)
	if err != nil {
		t.Fatal(err)
	}

	_ = saveContextAndHostData(opts.host, sessionID)

	_, token, _ := hostAndToken()

	if token != sessionID {
		t.Fatalf(errTempl, token, sessionID)
	}
}

func TestInvalidArgCombo(t *testing.T) {
	opts := &cfmLoginOptions{password: "yes", passwordStdin: true}

	err := validateArgs(opts)
	if err == nil {
		t.Fatal("should have got validation error")
	}

	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Error("wrong error returned")
	}
}

func TestValidArgCombo(t *testing.T) {
	opts := &cfmLoginOptions{password: "yes", passwordStdin: false}

	err := validateArgs(opts)
	if err != nil {
		t.Fatal("should have got validation error")
	}
}

func newTestServer(path string, h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewTLSServer(mux)
	mux.HandleFunc(path, h)

	return server
}
