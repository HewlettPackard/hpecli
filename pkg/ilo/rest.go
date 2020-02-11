// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/logger"
)

// restClient - generic REST api client
type restClient struct {
	Endpoint string
	Username string
	Password string
	APIKey   string
}

// NewClient - get a new network client
func newClient(host, user, password string) *restClient {
	return &restClient{Endpoint: host, Username: user, Password: password}
}

func (c *restClient) restAPICall(method, urlPath string, body io.Reader) ([]byte, error) {
	u, err := normalize(c.Endpoint + urlPath)
	if err != nil {
		logger.Critical("Unable to correctly parse the URL: %s", c.Endpoint+urlPath)
		return nil, err
	}

	logger.Debug("RestAPICall: %s - %s", method, u)

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	if strings.EqualFold(method, "post") {
		req.Header.Add("Content-Type", "application/json")
	}

	req.Header.Set("X-Auth-Token", c.APIKey)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// setup proxy
	proxyURL, err := http.ProxyFromEnvironment(req)
	if err != nil {
		return nil, fmt.Errorf("error with proxy: %v - %q", proxyURL, err)
	}

	if proxyURL != nil {
		tr.Proxy = http.ProxyURL(proxyURL)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		// login does post that responds with 201
		token := resp.Header.Get("X-Auth-Token")
		return []byte(token), nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote request failed with: \"%s\"", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dst := &bytes.Buffer{}
	if err := json.Indent(dst, data, "", "  "); err != nil {
		logger.Warning("Unable to pretty-print output.")
		return data, nil
	}

	return dst.Bytes(), nil
}

func normalize(u string) (string, error) {
	parsed, err := url.Parse(strings.ToLower(u))
	if err != nil {
		return "", err
	}

	if parsed.Path != "" {
		parsed.Path = strings.ReplaceAll(parsed.Path, "//", "/")
	}

	return parsed.String(), nil
}
