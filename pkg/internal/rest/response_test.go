// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package rest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const errTmpl = "Unexpected value - got=%s, want=%s"

func TestResponseType(t *testing.T) {
	ts := newTestServer("/", func(w http.ResponseWriter, r *http.Request) {
	})
	defer ts.Close()

	r, err := Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	got := reflect.TypeOf(r)
	want := reflect.TypeOf(&Response{})
	if got != want {
		t.Fatalf(errTmpl, got, want)
	}

}

func TestByteResponse(t *testing.T) {
	want := []byte("some string value abc123")
	ts := newTestServer("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(want)
	})
	defer ts.Close()

	r, err := Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	got := r.Bytes()
	if !bytes.Equal(got, want) {
		t.Fatalf(errTmpl, got, want)
	}
}

func TestJSONResponse(t *testing.T) {
	compactJSON := []byte(`{"a":"b","c":{"d":"e"}}`)
	want := []byte("{\n  " + `"a": "b",` + "\n" + `  "c": {` + "\n" + `    "d": "e"` + "\n  }\n}")
	ts := newTestServer("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(compactJSON)
	})
	defer ts.Close()

	r, err := Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	got := r.JSON()
	if !(bytes.Equal(got, []byte(want))) {
		t.Fatalf(errTmpl, got, want)
	}
}

func TestJSONPrettyPrintFailure(t *testing.T) {
	// bad formatted json causes pretty print failure
	want := []byte(`"a":"b","c":{"d":"e"}}`)
	ts := newTestServer("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(want)
	})
	defer ts.Close()

	r, err := Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	got := r.JSON()
	if !(bytes.Equal(got, want)) {
		t.Fatalf(errTmpl, got, want)
	}
}

func TestUnmarshall(t *testing.T) {
	compactJSON := []byte(`{"geo":"geovalue","token":"tokenvalue"}`)
	ts := newTestServer("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(compactJSON)
	})
	defer ts.Close()

	r, err := Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	type myType struct {
		Geo   string `json:"geo"`
		Token string `json:"token"`
	}

	var mt myType

	err = r.Unmarshall(&mt)
	if err != nil {
		t.Fatal(err)
	}

	if mt.Geo != "geovalue" {
		t.Fatal("unmarshall didn't get expected field value: geovalue")
	}
	if mt.Token != "tokenvalue" {
		t.Fatal("unmarshall didn't get expected field value: tokenvalue")
	}
}

func TestUnmarshallFails(t *testing.T) {
	//bad json causes unmarshall error
	compactJSON := []byte(`{`)
	ts := newTestServer("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(compactJSON)
	})
	defer ts.Close()

	r, err := Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	type myType struct {
		Geo   string `json:"geo"`
		Token string `json:"token"`
	}

	var mt myType

	err = r.Unmarshall(&mt)
	if err == nil {
		t.Fatal(err)
	}
}

func newTestServer(path string, h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc(path, h)

	return server
}
