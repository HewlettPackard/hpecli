// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	restHost     = "someHost"
	restUsername = "username"
	restPassword = "password"
	errTempl     = "got=%s, want=%s"
)

func TestNewClient(t *testing.T) {
	c := newClient(restHost, restUsername, restPassword)
	if c == nil {
		t.Fatal("expected client to not be nil")
	}

	if restHost != c.Endpoint {
		t.Fatal("restHost doesn't match")
	}

	if restUsername != c.Username {
		t.Fatal("restUsername doesn't match")
	}

	if restPassword != c.Password {
		t.Fatal("restPassword doesn't match")
	}
}

func TestBadHost(t *testing.T) {
	const badHost = ":badURL/missing/scheme"

	c := newClient(badHost, restUsername, restPassword)
	if _, err := c.restAPICall("POST", "", nil); err == nil {
		t.Fatalf("expected error parsing %s", badHost)
	}
}

func TestBadMethod(t *testing.T) {
	c := newClient(restHost, restUsername, restPassword)
	if _, err := c.restAPICall("\xb2", "/", nil); err == nil {
		t.Fatalf("expected error for bad method BAD")
	}
}

func TestGet(t *testing.T) {
	const want = "GOOD REQUEST"

	ts := newTestServer("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, want)
	})

	defer ts.Close()

	c := newClient(ts.URL, restUsername, restPassword)

	got, err := c.restAPICall("GET", "/", nil)
	if err != nil {
		t.Fatalf("Didn't get expected result data")
	}

	if string(got) != want {
		t.Fatalf(errTempl, got, want)
	}
}

func Test404Response(t *testing.T) {
	ts := newTestServer("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	defer ts.Close()

	c := newClient(ts.URL, restUsername, restPassword)

	_, err := c.restAPICall("GET", "/", nil)
	if err != nil && !strings.Contains(err.Error(), "404") {
		t.Fatal("didn't receive 404 error as expected")
	}
}

func TestEmptyBody(t *testing.T) {
	ts := newTestServer("/", func(w http.ResponseWriter, r *http.Request) {
		// not writing anything returns an empty 200 response
	})

	defer ts.Close()

	c := newClient(ts.URL, restUsername, restPassword)

	got, err := c.restAPICall("GET", "/", nil)
	if err != nil {
		t.Fatal("got unexpected error on empty response body")
	}

	want := make([]byte, 0)
	if !bytes.Equal(got, want) {
		t.Fatalf(errTempl, got, want)
	}
}

func TestNormalize(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "MixedCase",
			s:    "htTP://MixedCase/naMe",
			want: strings.ToLower("htTP://MixedCase/naMe"),
		},
		{
			name: "double slash",
			s:    "http://host//foo//bar//",
			want: "http://host/foo/bar/",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got, err := normalize(c.s)

			if err != nil {
				t.Fatalf("error not expected in normalizing: %s", c.s)
			}

			if got != c.want {
				t.Fatalf("normalize failed. "+errTempl, got, c.want)
			}
		})
	}
}

func newTestServer(path string, h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc(path, h)

	return server
}
