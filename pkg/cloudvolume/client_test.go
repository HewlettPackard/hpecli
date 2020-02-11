// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"fmt"
	"net/http"
	"testing"
)

const (
	clientHost     = "someClientHost"
	clientUsername = "username"
	clientPassword = "password"
)

func TestNewCVCClient(t *testing.T) {
	c := NewCVClient(clientHost, clientUsername, clientPassword)
	if c == nil {
		t.Fatal("expected client to not be nil")
	}

	if clientHost != c.Endpoint {
		t.Fatal("clientHost doesn't match")
	}

	if clientUsername != c.Username {
		t.Fatal("clientUsername doesn't match")
	}

	if clientPassword != c.Password {
		t.Fatal("clientPassword doesn't match")
	}
}

func TestMalformedResponseForLogin(t *testing.T) {
	const notJSON = "bad response"

	ts := newTestServer("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, notJSON)
	})

	defer ts.Close()

	c := NewCVClient(ts.URL, restUsername, restPassword)

	_, err := c.Login()
	if err == nil {
		t.Fatalf("Didn't get expected error on not json response")
	}
}

func TestTokeResponseForLogin(t *testing.T) {
	const want = "74dc0153-6daa-49ae-905e-cc59bff3225e"
	jsonResponse := fmt.Sprintf(`{"geo":"US", "token":"%s"}`, want)

	ts := newTestServer("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, jsonResponse)
	})

	defer ts.Close()

	c := NewCVClient(ts.URL, restUsername, restPassword)

	got, err := c.Login()
	if err != nil {
		t.Fatalf("unexpected error in login attempt")
	}

	if got != want {
		t.Fatalf(errTempl, got, want)
	}
}
