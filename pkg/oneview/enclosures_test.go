// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/store"
)

const encURL = "/rest/enclosures"

func TestAPIKeyPutInEnclosureRequest(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(encURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"sessionID":"%s"}`, sessionID)
		if r.Header.Get("Auth") != sessionID {
			t.Fatal("Expected to find \"auth\" header in request")
		}
	})

	defer server.Close()

	c := initContext(t)
	// set context to the test server host
	_ = c.SetAPIKey(server.URL, sessionID)

	// check is above in the http request handler side
	_ = getEnclosuresData()
}

func TestEnclosureClientRequestFails(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(encURL, func(w http.ResponseWriter, r *http.Request) {
		// cause the request to fail
		w.WriteHeader(http.StatusBadRequest)
	})

	defer server.Close()

	c := initContext(t)
	// set context to the test server host
	_ = c.SetAPIKey(server.URL, sessionID)

	// check is above in the http request handler side
	if err := getEnclosures(nil, nil); err == nil {
		t.Fatal("expected to get an error")
	}
}

func TestEnclosureJSONMarshallFails(t *testing.T) {
	const sessionID = "HERE_IS_A_ID"

	server := newTestServer(encURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// malformed json will cause the request to fail
		fmt.Fprint(w, `{"type":"server":["bad":"broken"]}`)
	})

	defer server.Close()

	c := initContext(t)
	// set context to the test server host
	_ = c.SetAPIKey(server.URL, sessionID)

	// check is above in the http request handler side
	if err := getEnclosuresData(); err == nil {
		t.Fatal("expected to get an error")
	}
}

func TestEnclosuresMissingAPIKey(t *testing.T) {
	// when the db is open, the get apikey will fail
	db, _ := store.Open()
	defer db.Close()

	err := getEnclosuresData()
	if err == nil {
		t.Fatal("should have retrieved error")
	}
}

func TestGetEnclosureByName(t *testing.T) {
	const name = "server-name"
	ovEnclosureData.name = name

	server := newTestServer(encURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"category":"enclosures","serialNumber":"SGH104X6J1","type":"EnclosureV7",`+
			`"uri":"/rest/enclosures/09SGH104X6J1","uuid":"09SGH104X6J1"}`)
	})

	defer server.Close()

	c := initContext(t)
	// set context to the test server host
	_ = c.SetAPIKey(server.URL, "sessionID")

	if err := getEnclosuresData(); err != nil {
		t.Fatal(err)
	}
}
