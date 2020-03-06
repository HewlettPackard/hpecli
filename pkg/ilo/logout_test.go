// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestLogoutHostPrefixAdded(t *testing.T) {
	// // clear everything from the mock store
	context.MockClear()

	mux := http.NewServeMux()
	server := httptest.NewTLSServer(mux)
	mux.HandleFunc("/rest/login-sessions", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	defer server.Close()

	host := strings.Replace(server.URL, "https://", "", 1)

	cmd := newLogoutCommand()
	cmd.SetArgs([]string{"--host", host})
	_ = cmd.Execute()

	_, err := getSessionData(server.URL)
	if !errors.Is(err, context.ErrorKeyNotFound) {
		t.Fatal("logout should delete the context")
	}
}

func TestLogoutSessionDataDeleted(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	pathURI := "/redfish/v1/sessionservice/sessions/fooSession"

	server := newTestServer(pathURI, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	defer server.Close()

	host := server.URL

	saveContextAndSessionData(&sessionData{server.URL, sessionID, server.URL + pathURI})

	err := runLogout(host)
	if err != nil {
		t.Fatal(err)
	}

	_, err = getSessionData(host)
	if err != context.ErrorKeyNotFound {
		t.Fatalf("expected ErrorKeyNotFound, but found %+v", err)
	}
}

func TestLogoutDefaultSessionDataDeleted(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	pathURI := "/redfish/v1/sessionservice/sessions/fooSession"

	server := newTestServer(pathURI, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	defer server.Close()

	host := ""

	saveContextAndSessionData(&sessionData{server.URL, sessionID, server.URL + pathURI})

	err := runLogout(host)
	if err != nil {
		t.Fatal(err)
	}

	_, err = getSessionData(host)
	if err != context.ErrorKeyNotFound {
		t.Fatalf("expected ErrorKeyNotFound, but found %+v", err)
	}
}
