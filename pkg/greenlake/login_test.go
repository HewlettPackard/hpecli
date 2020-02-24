// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestGlHostPrefixAdded(t *testing.T) {
	server := newTestServer("/identity/v1/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	glLoginData.host = strings.Replace(server.URL, "http://", "", 1)

	// this will fail with a remote call.. ignore the failure and
	// check the host string to ensure prefix addded
	_ = runGLLogin(nil, nil)

	if !strings.HasPrefix(glLoginData.host, "http://") {
		t.Fatalf("host should be prefixed with http scheme")
	}
}

func TestGLAccessTokenIsStored(t *testing.T) {
	const accessToken = "GreenLake_Access_Token"

	server := newTestServer("/identity/v1/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"access_token":"%s"}`, accessToken)
	})

	defer server.Close()

	glLoginData.host = server.URL

	err := runGLLogin(nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	c, _ := getData()
	if c.APIKey != accessToken {
		t.Fatalf(errTempl, c.APIKey, accessToken)
	}
}
