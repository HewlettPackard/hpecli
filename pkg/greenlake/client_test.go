// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

const (
	clientHost      = "greenlakeHost"
	clientID        = "username"
	clientSecretKey = "secretKey"
	clientGrantType = "grantType"
	clientToken     = "abcdefghij"
	clientTenantID  = "someTenantID"
	errTempl        = "got=%s, want=%s"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestNewGLCClient(t *testing.T) {
	c := NewGLClient(clientGrantType, clientID, clientSecretKey, clientTenantID, clientHost)
	if c == nil {
		t.Fatal("expected client to not be nil")
	}

	if clientGrantType != c.GrantType {
		t.Fatal("clientGrantType doesn't match")
	}

	if clientID != c.ClientID {
		t.Fatal("clientID doesn't match")
	}

	if clientSecretKey != c.ClientSecret {
		t.Fatal("clientPassword doesn't match")
	}

	if clientTenantID != c.TenantID {
		t.Fatal("clientTenantID doesn't match")
	}

	if clientHost != c.Host {
		t.Fatal("clientHost doesn't match")
	}
}

func TestNewGLClientFromAPIKey(t *testing.T) {
	c := NewGLClientFromAPIKey(clientHost, clientTenantID, clientToken)
	if c == nil {
		t.Fatal("expected client to not be nil")
	}

	if clientHost != c.Host {
		t.Fatal("clientHost doesn't match")
	}

	if clientTenantID != c.TenantID {
		t.Fatal("clientTenantID doesn't match")
	}

	if clientToken != c.APIKey {
		t.Fatal("clientToken doesn't match")
	}
}

func TestGLMalformedResponseForLogin(t *testing.T) {
	const notJSON = "bad response"

	ts := newTestServer("/identity/v1/token", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, notJSON)
	})

	defer ts.Close()

	c := NewGLClient(clientGrantType, clientID, clientSecretKey, clientTenantID, ts.URL)

	_, err := c.login()
	if err == nil {
		t.Fatalf("Didn't get expected error on not json response")
	}
}

func TestGLTokenResponseForLogin(t *testing.T) {
	const wantToken = "74dc0153-6daa-49ae-905e-cc59bff3225e-e"

	ts := newTestServer("/identity/v1/token", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"access_token":"%s"}`, wantToken)
	})

	defer ts.Close()

	c := NewGLClient(clientGrantType, clientID, clientSecretKey, clientTenantID, ts.URL)

	got, err := c.login()
	if err != nil {
		t.Fatalf("unexpected error in login attempt")
	}

	if got.Token != wantToken {
		t.Fatalf(errTempl, got.Token, wantToken)
	}
}

func TestGLAPIKeyInjected(t *testing.T) {
	const want = "dXNlcm5hbWU6dG9rZW4=="

	const userName = "userName"

	ts := newTestServer("/scim/v1/tenant/"+clientTenantID+"/Users", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `[{"userName":"%s"}`, userName)
	})

	defer ts.Close()
	c := NewGLClientFromAPIKey(ts.URL, clientTenantID, want)

	// checks are done on server side above
	_, _ = c.GetUsers()
}

func TestGetUsers(t *testing.T) {
	const compactJSON = `[{"active": false,"displayName": "anand.vuppuluri"}]`

	const want = "[\n  {\n    \"active\": false,\n    \"displayName\": \"anand.vuppuluri\"\n  }\n]"

	ts := newTestServer("/scim/v1/tenant/"+clientTenantID+"/Users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, compactJSON)
	})

	defer ts.Close()

	c := NewGLClientFromAPIKey(ts.URL, clientTenantID, clientToken)

	got, err := c.GetUsers()
	if err != nil {
		t.Fatalf("unexpected error in GetServiceRoot attempt")
	}

	s := string(got)
	// validate starts with with correct formatting
	if s != want {
		t.Fatalf(errTempl, got, want)
	}
}

func newTestServer(path string, h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc(path, h)

	return server
}
