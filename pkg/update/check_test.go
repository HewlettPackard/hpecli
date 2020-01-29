// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package update

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	contentType = "Content-Type"
	jsonType    = "application/json"
)

func TestIsUpdateAvailableEmptyLocalVersion(t *testing.T) {
	cases := []struct {
		name       string
		localVer   string
		remoteJSON string
		update     bool
	}{
		{
			name:       "no local version",
			localVer:   "0.0.0",
			remoteJSON: `{"version":"0.0.0"}`,
			update:     false,
		},
		{
			name:       "no local version",
			localVer:   "",
			remoteJSON: `{"version":"0.0.1"}`,
			update:     true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// erase cache for each test run
			cacheResponse = nil
			server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set(contentType, jsonType)
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, c.remoteJSON)
			})
			defer server.Close()

			got := IsUpdateAvailable()
			if got != c.update {
				t.Fatal("didn't get expected response")
			}
		})
	}
}

func TestIsUpdateAvailablInvalidURLErrors(t *testing.T) {
	// erase cache for each test run
	cacheResponse = nil
	versionURL = "://badScheme"

	got := IsUpdateAvailable()
	if got != false {
		t.Fatal("error in checkUpdate should generate false response")
	}
}

func TestCheckSkippedWithEnvSet(t *testing.T) {
	// erase cache for each test run
	cacheResponse = nil

	os.Setenv(EnvDisableUpdateCheck, "true")

	defer os.Unsetenv(EnvDisableUpdateCheck)

	got, _ := checkUpdate(&jsonSource{url: ""}, "")
	want := &CheckResponse{}

	//should return empty response because we skip everything
	//when the env var is set
	verifyCheckResponse(t, got, want)
}

func TestCheckUpdate(t *testing.T) {
	cases := []struct {
		name        string
		localVer    string
		remoteJSON  string
		errExpected bool
		want        *CheckResponse
	}{
		{
			name:       "remote greater than local",
			localVer:   "0.0.1",
			remoteJSON: `{"version":"0.1.0"}`,
			want: &CheckResponse{
				UpdateAvailable: true,
				RemoteVersion:   "0.1.0",
			},
		},
		{
			name:       "remote less than local",
			localVer:   "0.0.2",
			remoteJSON: `{"version":"0.0.1"}`,
			want: &CheckResponse{
				UpdateAvailable: false,
				RemoteVersion:   "0.0.1",
			},
		},
		{
			name:       "check all fields",
			localVer:   "0.1.2",
			remoteJSON: `{"version":"0.1.1","message":"update available","url":"https://foo.bar/update","publickey":"00001111","checksum":"120E0A"}`,
			want: &CheckResponse{
				UpdateAvailable: false,
				RemoteVersion:   "0.1.1",
				Message:         "update available",
				URL:             "https://foo.bar/update",
				PublicKey:       []byte{0x00, 0x00, 0x11, 0x11},
				CheckSum:        []byte{0x12, 0x0E, 0x0A},
			},
		},
		{
			name:        "missing local version",
			localVer:    "",
			errExpected: true,
		},
		{
			name:        "missing remote version",
			localVer:    "0.0.1",
			remoteJSON:  `{"message":"test will fail"}`,
			errExpected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// erase cache for each test run
			cacheResponse = nil
			server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set(contentType, jsonType)
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, c.remoteJSON)
			})
			defer server.Close()

			json := &jsonSource{
				url: versionURL,
			}

			got, err := checkUpdate(json, c.localVer)
			if err != nil {
				if c.errExpected {
					// got an error.. and expected an error
					return
				}
				// got an error but didn't expect it
				t.Fatal(err)
			}
			verifyCheckResponse(t, got, c.want)
		})
	}
}

func TestCachedCopyDoestRetrieveAgain(t *testing.T) {
	// erase cache for each test run
	cacheResponse = nil
	cc := 0

	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"version":"0.0.1"}`)
		cc++
		if cc > 1 {
			t.Fatal("Expected to only receive a single http request, but received 2")
		}
	})
	defer server.Close()

	// verify respnose is nil before initial get
	if cacheResponse != nil {
		t.Fatal("response to be unititialed before request")
	}

	got := IsUpdateAvailable()
	if got != true {
		t.Fatal("expected to see update available, but reported as not available")
	}
	// make sure response got populated
	if cacheResponse == nil {
		t.Fatal("response to be ititialed after request")
	}
	// save to check later
	r := cacheResponse

	got = IsUpdateAvailable()
	// just make sure it is still true
	if got != true {
		t.Fatal("expected to see update available, but reported as not available")
	}

	// see if it is the same cached copy
	if r != cacheResponse {
		t.Fatal("Expected to get the same copy on the second get request but didn't")
	}
}

func verifyCheckResponse(t *testing.T, got *CheckResponse, want *CheckResponse) {
	const tmpl = "got: %v, wanted: %v"

	if got.UpdateAvailable != want.UpdateAvailable {
		t.Fatalf(tmpl, got.UpdateAvailable, want.UpdateAvailable)
	}

	if got.RemoteVersion != want.RemoteVersion {
		t.Fatalf(tmpl, got.RemoteVersion, want.RemoteVersion)
	}

	if got.Message != want.Message {
		t.Fatalf(tmpl, got.Message, want.Message)
	}

	if got.URL != want.URL {
		t.Fatalf(tmpl, got.URL, want.URL)
	}

	if bytes.Equal(got.PublicKey, want.PublicKey) {
		t.Fatalf(tmpl, got.PublicKey, want.PublicKey)
	}

	if !bytes.Equal(got.CheckSum, want.CheckSum) {
		t.Fatalf(tmpl, got.CheckSum, want.CheckSum)
	}
}

func newTestServer(h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	versionURL = fmt.Sprintf("%s%s", server.URL, versionPath)
	mux.HandleFunc(versionPath, h)

	return server
}
