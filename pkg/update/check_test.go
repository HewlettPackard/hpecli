// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package update

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tcnksm/go-latest"
)

func TestCheckUpdate(t *testing.T) {

	cases := []struct {
		name      string
		localVer  string
		remoteVer string
		want      *latest.CheckResponse
	}{
		{
			name:      "remote greater than local",
			localVer:  "0.0.1",
			remoteVer: `{"version":"0.1.0"}`,
			want: &latest.CheckResponse{
				Current:  "0.1.0",
				Outdated: true,
				Latest:   false,
				New:      false,
			},
		},
		{
			name:      "remote less than local",
			localVer:  "0.0.2",
			remoteVer: `{"version":"0.0.1"}`,
			want: &latest.CheckResponse{
				Current:  "0.0.1",
				Outdated: false,
				Latest:   true,
				New:      true,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			response = nil
			server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, c.remoteVer)
			})
			defer server.Close()

			json := &latest.JSON{
				URL: versionURL,
			}

			got, err := checkUpdate(json, c.localVer)
			if err != nil {
				t.Fatal(err)
			}
			validate(t, got, c.want)
		})
	}
}

func TestCachedCopyDoestRetrieveAgain(t *testing.T) {
	response = nil
	cc := 0

	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"version":"0.0.1"}`)
		cc++
		if cc > 1 {
			t.Fatal("Expected to only recieve a single http request, but recieved 2")
		}
	})
	defer server.Close()

	// verify respnose is nil before initial get
	if response != nil {
		t.Fatal("response to be unititialed before request")
	}

	got := IsUpdateAvailable()
	if got != true {
		t.Fatal("expected to see update available, but reported as not available")
	}
	// make sure response got populated
	if response == nil {
		t.Fatal("response to be ititialed after request")
	}
	// save to check later
	r := response

	got = IsUpdateAvailable()

	// see if it is the same cached copy
	if r != response {
		t.Fatal("Expected to get the same copy on the second get request but didn't")
	}
}

func TestErrorReturnsFalse(t *testing.T) {
	response = nil
	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, `{"version":"0.0.1"}`)
	})

	got := IsUpdateAvailable()
	// server responding with 404 should reply "false" - in that it can't
	// detect an updated version
	if got != false {
		t.Fatal("Request error should have returned update not available")
	}
	server.Close()
}

func validate(t *testing.T, got *latest.CheckResponse, want *latest.CheckResponse) {
	if got.Current != want.Current {
		t.Fatal(fmt.Sprintf("got: %v, wanted: %v", got.Current, want.Current))
	}
	if got.Outdated != want.Outdated {
		t.Fatal(fmt.Sprintf("got: %v, wanted: %v", got.Outdated, want.Outdated))
	}
	if got.Latest != want.Latest {
		t.Fatal(fmt.Sprintf("got: %v, wanted: %v", got.Latest, want.Latest))
	}
	if got.New != want.New {
		t.Fatal(fmt.Sprintf("got: %v, wanted: %v", got.New, want.New))
	}
}

func newTestServer(h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	versionURL = fmt.Sprintf("%s%s", server.URL, versionPath)
	mux.HandleFunc(versionPath, h)
	return server
}
