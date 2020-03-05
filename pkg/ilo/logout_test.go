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

func TestLogoutHostPrefixAdded(t *testing.T) {
	host := "127.0.0.1"

	// this will fail with a remote call.. ignore the failure and
	// check the host string to ensure prefix addded
	_ = runILOLogout(&host)

	if !strings.HasPrefix(host, "https://") {
		t.Fatalf("host should be prefixed with http scheme")
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

	err := runILOLogout(&host)
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

	err := runILOLogout(&host)
	if err != nil {
		t.Fatal(err)
	}

	_, err = getSessionData(host)
	if err != context.ErrorKeyNotFound {
		t.Fatalf("expected ErrorKeyNotFound, but found %+v", err)
	}
}
