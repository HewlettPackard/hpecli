// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

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

func TestHostPrefixAdded(t *testing.T) {
	// clear everything from the mock store
	context.MockClear()

	mux := http.NewServeMux()
	server := httptest.NewTLSServer(mux)
	mux.HandleFunc("/rest/login-sessions", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"sessionID":"%s"}`, "sessionID")
	})

	defer server.Close()

	host := strings.Replace(server.URL, "https://", "", 1)

	cmd := newLoginCommand()
	cmd.SetArgs([]string{"--host", host, "--username", "admin", "--password", "somePse"})
	_ = cmd.Execute()

	// check the db to make sure it was persisted
	got, err := getContext()
	if err != nil {
		t.Fatal(err)
	}

	if got != "https://"+host {
		t.Fatal("context didn't get set correctly")
	}
}

func TestAPIKeyIsStored(t *testing.T) {
	const sessionID = "SomeSessionID"

	server := newTestServer("/rest/login-sessions", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"sessionID":"%s"}`, sessionID)
	})

	defer server.Close()

	opts := &ovLoginOptions{
		host:     server.URL,
		password: "blah blah",
	}

	err := runLogin(opts)
	if err != nil {
		t.Fatal(err)
	}

	_, token, _ := hostAndToken()

	if token != sessionID {
		t.Fatalf(errTempl, token, sessionID)
	}
}

func TestInvalidArgCombo(t *testing.T) {
	opts := &ovLoginOptions{password: "yes", passwordStdin: true}

	err := validateArgs(opts)
	if err == nil {
		t.Fatal("should have got validation error")
	}

	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Error("wrong error returned")
	}
}

func newTestServer(path string, h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc(path, h)

	return server
}
