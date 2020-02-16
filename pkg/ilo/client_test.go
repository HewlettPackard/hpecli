// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	clientHost     = "someClientHost"
	clientUsername = "username"
	clientPassword = "password"
	clientToken    = "ljwer;lkjl23j4lk3l;jlk"
	errTempl       = "got=%s, want=%s"
)

func TestNewILOCClient(t *testing.T) {
	c := NewILOClient(clientHost, clientUsername, clientPassword)
	if c == nil {
		t.Fatal("expected client to not be nil")
	}

	if clientHost != c.Host {
		t.Fatal("clientHost doesn't match")
	}

	if clientUsername != c.Username {
		t.Fatal("clientUsername doesn't match")
	}

	if clientPassword != c.Password {
		t.Fatal("clientPassword doesn't match")
	}
}

func TestNewILOClientFromAPIKey(t *testing.T) {
	c := NewILOClientFromAPIKey(clientHost, clientToken)
	if c == nil {
		t.Fatal("expected client to not be nil")
	}

	if clientHost != c.Host {
		t.Fatal("clientHost doesn't match")
	}

	if clientToken != c.APIKey {
		t.Fatal("clientToken doesn't match")
	}
}

func TestMalformedResponseForLogin(t *testing.T) {
	const notJSON = "bad response"

	ts := newTestServer("/redfish/v1/sessionservice/sessions/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, notJSON)
	})

	defer ts.Close()

	c := NewILOClient(ts.URL, clientUsername, clientPassword)

	_, err := c.Login()
	if err == nil {
		t.Fatalf("Didn't get expected error on not json response")
	}
}

func TestTokenResponseForLogin(t *testing.T) {
	const want = "74dc0153-6daa-49ae-905e-cc59bff3225e"

	ts := newTestServer("/redfish/v1/sessionservice/sessions/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("x-auth-token", want)
		w.WriteHeader(http.StatusCreated)
	})

	defer ts.Close()

	c := NewILOClient(ts.URL, clientUsername, clientPassword)

	got, err := c.Login()
	if err != nil {
		t.Fatalf("unexpected error in login attempt")
	}

	if got != want {
		t.Fatalf(errTempl, got, want)
	}
}

func TestAPIKeyInjected(t *testing.T) {
	const want = "dXNlcm5hbWU6dG9rZW4="

	ts := newTestServer("/redfish/v1/", func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("X-Auth-Token")
		if got != want {
			t.Fatal("didn't get expected auth header")
		}
	})

	defer ts.Close()

	c := NewILOClientFromAPIKey(ts.URL, want)

	// checks are done on server side above
	_, _ = c.GetServiceRoot()
}

func TestGetServiceRoot(t *testing.T) {
	const compactJSON = `{"@odata.context":"/redfish/v1/$metadata#ServiceRoot.ServiceRoot","@odata.id":"/redfish/v1/",` +
		`"@odata.type":"#ServiceRoot.1.0.0.ServiceRoot","AccountService":{"@odata.id":"/redfish/v1/AccountService/"},` +
		`"Chassis":{"@odata.id":"/redfish/v1/Chassis/"}}`

	const want = "{\n  \"@odata.context\": \"/redfish/v1/$metadata#ServiceRoot.ServiceRoot\",\n"

	ts := newTestServer("/redfish/v1/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, compactJSON)
	})

	defer ts.Close()

	c := NewILOClientFromAPIKey(ts.URL, "someAPIKey")

	got, err := c.GetServiceRoot()
	if err != nil {
		t.Fatalf("unexpected error in GetServiceRoot attempt")
	}

	s := string(got)
	// validate starts with with correct formatting
	if !strings.HasPrefix(s, want) {
		t.Fatalf(errTempl, got, want)
	}
}

func newTestServer(path string, h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc(path, h)

	return server
}
