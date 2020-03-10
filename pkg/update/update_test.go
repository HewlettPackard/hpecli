// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package update

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testServerURL = "test-server"

func TestUpdateWithNoUpdateAvailable(t *testing.T) {
	cacheResponse = nil
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	mux.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf(`{"version":"0.0.0","url":"%s/download"}`, ts.URL))
	})

	defer ts.Close()

	versionURL = fmt.Sprintf("%s%s", ts.URL, "/json")

	if err := runUpdate("0.0.0"); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateWithErrorInCheckUpdate(t *testing.T) {
	cacheResponse = nil
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	mux.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf(`{"version":"0.0.0","url":"%s/download"}`, ts.URL))
	})

	defer ts.Close()

	cmd := NewUpdateCommand("A")
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected failure")
	}
}

func TestDownloadUpdate(t *testing.T) {
	cases := []struct {
		name        string
		param       *CheckResponse
		want        string
		errExpected bool
	}{
		{
			name: "download works",
			param: &CheckResponse{
				URL: testServerURL,
			},
			want:        "value doesn't matter as long as it matches",
			errExpected: false,
		},
		{
			name: "expected error in response body",
			param: &CheckResponse{
				URL: "",
			},
			errExpected: true,
		},
		{
			name: "expected error in incorrect checksum",
			param: &CheckResponse{
				URL:      testServerURL,
				CheckSum: []byte{0x00},
			},
			errExpected: true,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, c.want)
			}))

			defer ts.Close()

			if c.param.URL == testServerURL {
				c.param.URL = ts.URL
			}

			err := downloadUpdate(c.param)
			if err != nil {
				if c.errExpected {
					return
				}
				t.Fatal(err)
			}
		})
	}
}

//nolint:gocritic  // complex test method NOSONAR
func TestGetResponseBody(t *testing.T) {
	cases := []struct {
		name        string
		url         string
		want        string
		errExpected bool
	}{
		{
			name:        "read body works",
			url:         testServerURL,
			want:        "value doesn't matter as long as it matches",
			errExpected: false,
		},
		{
			name:        "bad url",
			url:         "://missing/scheme",
			errExpected: true,
		},
		{
			name:        "empty url",
			url:         "",
			errExpected: true,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, c.want)
			}))
			defer ts.Close()

			if c.url == testServerURL {
				c.url = ts.URL
			}

			body, err := getResponseBody(c.url)
			if err != nil {
				if c.errExpected {
					return
				}
				t.Fatal(err)
			}

			buf, err := ioutil.ReadAll(body)
			if err != nil {
				t.Fatal(err)
			}
			got := string(buf)
			if got != c.want {
				t.Fatalf("Didn't get expected response.  got=%v, want=%v", got, c.want)
			}
		})
	}
}
