// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestHostPrefixAdded(t *testing.T) {
	// clear everything from the mock store
	context.MockClear()

	mux := http.NewServeMux()
	server := httptest.NewTLSServer(mux)
	mux.HandleFunc("/redfish/v1/sessionservice/sessions/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("x-auth-token", "someToken")
		w.WriteHeader(http.StatusCreated)
	})

	defer server.Close()

	host := strings.Replace(server.URL, "https://", "", 1)

	cmd := newLoginCommand()
	cmd.SetArgs([]string{"--host", host, "-u", "user", "-p", "pass"})
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
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer("/redfish/v1/sessionservice/sessions/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("x-auth-token", sessionID)
		w.WriteHeader(http.StatusCreated)
	})

	defer server.Close()

	opts := iloLoginOptions{
		host:     server.URL,
		password: "blah",
	}

	err := runLogin(&opts)
	if err != nil {
		t.Fatal(err)
	}

	d, _ := defaultSessionData()
	if d.Host != opts.host {
		t.Fatalf(errTempl, d.Host, opts.host)
	}

	if d.Token != sessionID {
		t.Fatalf(errTempl, d.Token, sessionID)
	}
}

func TestInvalidArgCombo(t *testing.T) {
	opts := &iloLoginOptions{password: "yes", passwordStdin: true}

	err := validateArgs(opts)
	if err == nil {
		t.Fatal("should have got validation error")
	}

	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Error("wrong error returned")
	}
}
