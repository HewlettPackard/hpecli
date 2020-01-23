// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-version"
)

type jsonSource struct {
	url string
}

type jsonResponse struct {
	Version   string `json:"version"`
	Message   string `json:"message"`
	UpdateURL string `json:"url"`
	PublicKey string `json:"publickey"`
	CheckSum  string `json:"checksum"`
}

func (j *jsonSource) validate() error {
	if j.url == "" {
		return fmt.Errorf("Remote URL must be set")
	}
	if _, err := url.Parse(j.url); err != nil {
		return fmt.Errorf("%s is invalid URL: %s", j.url, err.Error())
	}
	return nil
}

func (j *jsonSource) get() (*remoteResponse, error) {
	client := client()

	req, err := request(j.url)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	jsonResponse, err := decodeResponse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to decode json response: %v", err)
	}
	result, err := mapResult(jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("unable to map json to remoteResponse: %v", err)
	}

	return result, nil
}

func client() *http.Client {
	const defaultTimeout = 15 * time.Second
	// Create client that will use env proxy settings
	// and ha a 15 second timeout
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
		Timeout: defaultTimeout,
	}
	return client
}

func request(url string) (*http.Request, error) {
	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	return req, nil
}

func decodeResponse(r io.Reader) (*jsonResponse, error) {
	result := &jsonResponse{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func mapResult(j *jsonResponse) (*remoteResponse, error) {
	v, err := version.NewVersion(j.Version)
	if err != nil {
		return nil, fmt.Errorf("Didn't get a remote version. %v", err)
	}
	return &remoteResponse{
		version:   v,
		message:   j.Message,
		updateURL: j.UpdateURL,
		publicKey: []byte(j.PublicKey),
		checkSum:  j.CheckSum,
	}, nil
}
