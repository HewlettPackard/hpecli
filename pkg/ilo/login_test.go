// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"net/http"
	"strings"
	"testing"
)

func TestHostPrefixAdded(t *testing.T) {
	server := newTestServer("/redfish/v1/sessionservice/sessions/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	iloLoginData.host = strings.Replace(server.URL, "http://", "", 1)

	// this will fail with a remote call.. ignore the failure and
	// check the host string to ensure prefix addded
	_ = runILOLogin(nil, nil)

	if !strings.HasPrefix(iloLoginData.host, "https://") {
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

	iloLoginData.host = server.URL

	err := runILOLogin(nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	c := iloContext()

	// sessionId is stored - so get it and verify it
	_, gotAPIKey, _ := c.APIKey()
	if gotAPIKey != sessionID {
		t.Fatalf(errTempl, gotAPIKey, sessionID)
	}
}
