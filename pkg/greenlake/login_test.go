// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

const uriPath = "/identity/v1/token"

func TestGLHostPrefixAddedForLogin(t *testing.T) {
	// clear everything from the mock store
	context.MockClear()

	server := newTestServer(uriPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"access_token":"accessToken"}`)
	})

	defer server.Close()

	host := strings.Replace(server.URL, "http://", "", 1)

	cmd := newLoginCommand()
	cmd.SetArgs([]string{"--host", host, "--tenantid", "id", "--userid", "user", "--secretkey", "key"})
	_ = cmd.Execute()

	// check to ensure context value gets the http scheme added
	got, _ := getContext()
	if got != "http://"+host {
		t.Error("context value didn't get http scheme prefix")
	}
}

func TestGLAccessTokenIsStored(t *testing.T) {
	const accessToken = "GreenLake_Access_Token"

	server := newTestServer(uriPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"access_token":"%s"}`, accessToken)
	})

	defer server.Close()

	cmd := newLoginCommand()
	cmd.SetArgs([]string{"--host", server.URL, "--tenantid", "id", "--userid", "user", "--secretkey", "key"})

	_ = cmd.Execute()

	d, _ := defaultSessionData()
	if d.Host != server.URL {
		t.Fatalf(errTempl, d.Host, server.URL)
	}

	if d.Token != accessToken {
		t.Fatalf(errTempl, d.Token, accessToken)
	}
}

func TestHTTPFailure(t *testing.T) {
	server := newTestServer(uriPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	defer server.Close()

	opts := &glLoginOptions{
		secretKey: "key",
	}

	err := runLogin(opts)
	if err == nil {
		t.Fatal("failed http request should fail the login request")
	}
}

func TestInvalidArgCombo(t *testing.T) {
	opts := &glLoginOptions{secretKey: "yes", secretKeyStdin: true}

	err := validateArgs(opts)
	if err == nil {
		t.Fatal("should have got validation error")
	}

	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Error("wrong error returned")
	}
}
