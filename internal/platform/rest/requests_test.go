// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package rest

import (
	"net/http"
	"testing"
)

func TestRequests(t *testing.T) {
	ts := newTestServer("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	})
	defer ts.Close()

	r, err := Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != http.StatusBadGateway {
		t.Fatal(err)
	}
}

func TestOptionsExecuted(t *testing.T) {
	fn1 := func(r *Request) {
		r.SetBasicAuth("username", "password")
	}

	fn2 := AddJSONMimeType()

	fn3 := AddHeaders("someKey", "someValue")

	fn4 := AllowSelfSignedCerts()

	ts := newTestServer("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	defer ts.Close()

	r, err := Get(ts.URL, fn1, fn2, fn3, fn4)
	if err != nil {
		t.Fatal(err)
	}

	u, p, ok := r.Request.BasicAuth()
	if !ok || u != "username" || p != "password" {
		t.Fatal("didn't get basic auth values after set")
	}

	if r.Request.Header.Get("content-type") != "application/json" {
		t.Fatal("didn't get header content-type value after set")
	}

	if r.Request.Header.Get("someKey") != "someValue" {
		t.Fatal("didn't get header someKey  value after set")
	}
}

func TestPost(t *testing.T) {
	ts := newTestServer("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	defer ts.Close()

	r, err := Post(ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != http.StatusCreated {
		t.Fatalf("expected status: %v", http.StatusCreated)
	}
}

func TestDelete(t *testing.T) {
	ts := newTestServer("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	defer ts.Close()

	r, err := Delete(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != http.StatusOK {
		t.Fatalf("expected status: %v", http.StatusOK)
	}
}
