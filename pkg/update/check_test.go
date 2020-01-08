// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package update

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tcnksm/go-latest"
)

const v1 = "v0.0.1"

func TestIsUpdateAvailable(t *testing.T) {

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
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			versionURL = fmt.Sprintf("%s%s", server.URL, versionPath)

			mux.HandleFunc(versionPath, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, c.remoteVer)
			})

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
