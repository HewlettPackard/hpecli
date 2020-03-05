// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"net/http"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestHostPrefixAdded(t *testing.T) {
	server := newTestServer("/redfish/v1/sessionservice/sessions/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	opts := iloLoginOptions{
		host:     strings.Replace(server.URL, "http://", "", 1),
		password: "blah",
	}

	// this will fail with a remote call.. ignore the failure and
	// check the host string to ensure prefix addded
	_ = runLogin(&opts)

	if !strings.HasPrefix(opts.host, "https://") {
		t.Fatalf("host should be prefixed with http scheme")
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

func TestTwoPasswordOptions(t *testing.T) {
	opts := &iloLoginOptions{
		password:      "pswd",
		passwordStdin: true,
	}

	err := handlePasswordOptions(opts)
	if err == nil {
		t.Error("expected error")
	}

	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Error("unexpected error text")
	}
}
