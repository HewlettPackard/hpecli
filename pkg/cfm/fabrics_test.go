package cfm

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// TestNewFabricsCommand tests the command to get fabrics
func TestGewFabricsCommand(t *testing.T) {
	cmd := getFabricsCommand()
	if cmd.Use != "fabrics" {
		t.Errorf("command name appears to be wrong")
	}
}

// TestGetFabrics ...
func TestGetFabrics(t *testing.T) {

	server := newTestServer("/api/v1/fabrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader((http.StatusOK))
		fmt.Fprintf(w, `{"result": []}`)
	})

	err := saveContextAndHostData(strings.ReplaceAll(server.URL, "https://", ""), "API Token")
	if err != nil {
		t.Errorf("Could not set the context")
	}

	err = getFabrics()

	if err != nil {
		t.Errorf("getFabrics failed")
	}

	server = newTestServer("/api/v1/fabrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader((http.StatusInternalServerError))
		fmt.Fprintf(w, `{"result": "Internal Server Error"}`)
	})

	err = saveContextAndHostData(strings.ReplaceAll(server.URL, "https://", ""), "API Token")
	if err != nil {
		t.Errorf("Could not set the context")
	}

	err = getFabrics()

	if err == nil {
		t.Errorf("GetFabrics failed")
	}

	err = deleteSavedHostData(strings.ReplaceAll(server.URL, "https://", ""))

	err = getFabrics()
	if err == nil {
		t.Errorf("getfabrics failed")
	}
}
