// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/context"
	"github.com/HewlettPackard/hpecli/pkg/db"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestRunCVGetVolumesWithMissingAPIKey(t *testing.T) {
	f := db.KeystoreLocation()
	// delete the default store
	_ = os.Remove(f)

	err := runCVGetVolumes(nil, nil)
	if err == nil {
		t.Fatal(err)
	}
}

func TestRunCVGetVolumes(t *testing.T) {
	const jsonResponse = `{"data":[{"attributes":{"app_uuid":"","assigned_initiators":[],"cloud_accounts"` +
		`:[{"href":"https://demo.cloudvolumes.hpe.com/api/v2/session/cloud_accounts/F2aSi8dNerU5T3zAatrNc` +
		`8z2cevknLBWYSnyOrgg","id":"F2aSi8dNerU5T3zAatrNc8z2cevknLBWYSnyOrgg"}]}}]}`

	ts := newTestServer("/api/v2/cloud_volumes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, jsonResponse)
	})

	defer ts.Close()

	f := db.KeystoreLocation()
	// delete the default store
	_ = os.Remove(f)

	storeContext(ts.URL, "someAPIKey")

	err := runCVGetVolumes(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}
