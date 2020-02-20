// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/context"
)

const (
	clientHost     = "someClientHost"
	clientUsername = "username"
	clientPassword = "password"
	clientToken    = "lkjsdfjka;sdfjlasdjkf"
	errTempl       = "got=%s, want=%s"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func newTestServer(path string, h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc(path, h)

	return server
}

func TestNewCVCClient(t *testing.T) {
	c := NewCVClient(clientHost, clientUsername, clientPassword)
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

func TestNewCVClientFromAPIKey(t *testing.T) {
	c := NewCVClientFromAPIKey(clientHost, clientToken)
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

	ts := newTestServer("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, notJSON)
	})

	defer ts.Close()

	c := NewCVClient(ts.URL, clientUsername, clientPassword)

	_, err := c.Login()
	if err == nil {
		t.Fatalf("Didn't get expected error on not json response")
	}
}

func TestTokenResponseForLogin(t *testing.T) {
	const want = "74dc0153-6daa-49ae-905e-cc59bff3225e"
	jsonResponse := fmt.Sprintf(`{"geo":"US", "token":"%s"}`, want)

	ts := newTestServer("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, jsonResponse)
	})

	defer ts.Close()

	c := NewCVClient(ts.URL, clientUsername, clientPassword)

	got, err := c.Login()
	if err != nil {
		t.Fatalf("unexpected error in login attempt")
	}

	if got != want {
		t.Fatalf(errTempl, got, want)
	}
}

func TestAPIKeyInjected(t *testing.T) {
	const apiKey = "dXNlcm5hbWU6dG9rZW4="
	// header value is base64 encoding of "username:dXNlcm5hbWU6dG9rZW4="
	want := "Basic dXNlcm5hbWU6ZFhObGNtNWhiV1U2ZEc5clpXND0="

	ts := newTestServer("/api/v2/cloud_volumes", func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("Authorization")
		if got != want {
			t.Fatal("didn't get expected auth header")
		}
	})

	defer ts.Close()

	c := NewCVClientFromAPIKey(ts.URL, apiKey)

	// checks are done on server side above
	_, _ = c.GetCloudVolumes()
}

func TestGetCloudVolumes(t *testing.T) {
	const compactJSON = `{"data":[{"cloud_accounts":[{"href":"https://demo.cloudvolumes.hpe.com/api/v2/` +
		`session/cloud_accounts/F2aSi8dNerU5T3zAatrNc8z2cevknLBWYSnyOrgg","id":` +
		`"F2aSi8dNerU5T3zAatrNc8z2cevknLBWYSnyOrgg"}]}]}`

	const want = "{\n  \"data\": ["

	ts := newTestServer("/api/v2/cloud_volumes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, compactJSON)
	})

	defer ts.Close()

	c := NewCVClientFromAPIKey(ts.URL, "someAPIKey")

	got, err := c.GetCloudVolumes()
	if err != nil {
		t.Fatalf("unexpected error in getCloudVolumes attempt")
	}

	s := string(got)
	// validate starts with with correct formatting
	if !strings.HasPrefix(s, want) {
		t.Fatalf(errTempl, got, want)
	}
}
