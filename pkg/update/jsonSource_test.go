// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package update

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-version"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name        string
		source      source
		errExpected bool
		errMsg      string
	}{
		{
			name:        "emptyURL",
			source:      &jsonSource{url: ""},
			errExpected: true,
			errMsg:      "expected error on empty URL",
		},
		{
			name:        "invalid scheme",
			source:      &jsonSource{url: "://invalid\\url"},
			errExpected: true,
			errMsg:      "expected error on invalid url",
		},
		{
			name:        "valid URL doesn't error",
			source:      &jsonSource{url: "http://hpe.com"},
			errExpected: false,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			err := c.source.validate()
			if err == nil && c.errExpected {
				t.Fatal(c.errMsg)
			}
		})
	}
}

func TestGetMalformedURL(t *testing.T) {
	s := &jsonSource{url: "://bad.url"}

	_, err := s.get()
	if err == nil {
		t.Fatal("expected err with bad url")
	}
}

func TestGetNildURL(t *testing.T) {
	s := &jsonSource{}

	_, err := s.get()
	if err == nil {
		t.Fatal("expected err with bad url")
	}
}

func TestGetStatusNotOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	s := &jsonSource{url: ts.URL}

	_, err := s.get()
	if err == nil {
		t.Fatal("error expected, but didn't get one")
	}
}

func TestGetDecodeErrorWithMalformedJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"version":`)
	}))
	defer ts.Close()

	s := &jsonSource{url: ts.URL}

	_, err := s.get()
	if err == nil {
		t.Fatal("expected err with bad json response")
	}
}

func TestGetWithoutVersionInResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"message":"Some Message","url":"http://some.url","publickey":"00:11:22:33",`+
			`"checksum":"120EA8A25E5D487BF68B5F7096440019"}`)
	}))
	defer ts.Close()

	s := &jsonSource{url: ts.URL}

	_, err := s.get()
	if err == nil {
		t.Fatal("expected err with missing version")
	}
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"version":"0.1.1","message":"Some Message","url":"http://another.url","publickey":"00:11:22:33",`+
			`"checksum":"120EA8A25E5D487BF68B5F7096440019"}`)
	}))
	defer ts.Close()

	s := &jsonSource{url: ts.URL}

	resp, err := s.get()
	if err != nil {
		t.Fatal("error not expected")
	}

	// assume decode works.. just check a couple of fields
	if resp.version.String() != "0.1.1" {
		t.Fatal("version didn't decode as expected")
	}

	if resp.updateURL != "http://another.url" {
		t.Fatal("updateURL didn't decode as expected")
	}
}

//nolint:gocognit NOSONAR // long test method
func TestMapResult(t *testing.T) {
	cases := []struct {
		name        string
		json        *jsonResponse
		version     *version.Version
		message     string
		updateURL   string
		publicKey   []byte
		checkSum    []byte
		errExpected bool
	}{
		{
			name:        "empty version",
			json:        &jsonResponse{Version: ""},
			errExpected: true,
		},
		{
			name: "fields match",
			json: &jsonResponse{Version: "0.1.1", Message: "Some Message", UpdateURL: "http://github.com/update",
				PublicKey: "00112233", CheckSum: "120E0A8A25E5"},
			version:     ver("0.1.1"),
			message:     "Some Message",
			updateURL:   "http://github.com/update",
			publicKey:   []byte{0x00, 0x11, 0x22, 0x33},
			checkSum:    []byte{0x12, 0x0E, 0x0A, 0x8A, 0x25, 0xE5},
			errExpected: false,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got, err := mapResult(c.json)
			if err == nil && c.errExpected {
				t.Fatal("Error expected but not found")
			}
			if err != nil && !c.errExpected {
				t.Fatal("Error found and not expected")
			}
			if err != nil && c.errExpected {
				// got an error and expected to get one
				return
			}
			if !got.version.Equal(c.version) {
				t.Fatalf("Version doesn't match.  got=%v, want=%v", got.version, c.version)
			}
			if got.message != c.message {
				t.Fatalf("Message doesn't match.  got=%v, want=%v", got.message, c.message)
			}
			if got.updateURL != c.updateURL {
				t.Fatalf("updateURL doesn't match.  got=%v, want=%v", got.updateURL, c.updateURL)
			}
			if !bytes.Equal(got.publicKey, c.publicKey) {
				t.Fatalf("publicKey doesn't match.  got=%v, want=%v", got.publicKey, c.publicKey)
			}
			if !bytes.Equal(got.checkSum, c.checkSum) {
				t.Fatalf("checkSum doesn't match.  got=%v, want=%v", got.checkSum, c.checkSum)
			}
		})
	}
}

func ver(v string) *version.Version {
	r, _ := version.NewVersion(v)
	return r
}

func TestDecodeField(t *testing.T) {
	if decodeField("key", "ZZZ") != nil {
		t.Fatal("expected nil on failure")
	}
}
