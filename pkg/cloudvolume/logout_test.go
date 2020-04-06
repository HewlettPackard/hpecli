// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"errors"
	"net/http"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestLogoutHost(t *testing.T) {

	cmd := newLogoutCommand()
	cmd.SetArgs([]string{"--host", cvDefaultHost})
	_ = cmd.Execute()

	_, err := hostData(cvDefaultHost)
	if !errors.Is(err, context.ErrorKeyNotFound) {
		t.Fatal("logout should delete the context")
	}
}

func TestLogoutNoHost(t *testing.T) {

	cmd := newLogoutCommand()
	_ = cmd.Execute()

	_, err := hostData(cvDefaultHost)
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

	saveData(host, sessionID)

	err := runLogout(host)
	if err != nil {
		t.Fatal(err)
	}

	_, err = hostData(host)
	if err != context.ErrorKeyNotFound {
		t.Fatalf("expected ErrorKeyNotFound, but found %+v", err)
	}
}

