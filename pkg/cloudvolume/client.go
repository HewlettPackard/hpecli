// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/rest"
)

type CVClient struct {
	Host     string
	Username string
	Password string
	APIKey   string
	*rest.Request
}

func newCVClient(host, username, password string) *CVClient {
	return &CVClient{
		Host:     host,
		Username: username,
		Password: password,
	}
}

func newCVClientFromAPIKey(host, token string) *CVClient {
	return &CVClient{
		Host:   host,
		APIKey: token,
	}
}

func (c *CVClient) login() (string, error) {
	const uriPath = "/auth/login"

	postBody := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, c.Username, c.Password)

	resp, err := rest.Post(c.Host+uriPath, strings.NewReader(postBody), rest.AddJSONMimeType())
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to create login sessions to greenlake.  Repsponse was: %+v", resp.Status)
	}

	type loginResponse struct {
		Geo   string `json:"geo"`
		Token string `json:"token"`
	}

	var lr loginResponse

	err = resp.Unmarshall(&lr)
	if err == nil {
		return lr.Token, nil
	}

	return "", fmt.Errorf("unable to get response from login command")
}

func (c *CVClient) getCloudVolumes() ([]byte, error) {
	const uriPath = "/api/v2/cloud_volumes"

	resp, err := rest.Get(c.Host+uriPath, c.addAuth())
	if err != nil {
		return []byte{}, err
	}

	return resp.JSON(), nil
}

func (c *CVClient) addAuth() func(*rest.Request) {
	return func(r *rest.Request) {
		r.Request.SetBasicAuth("username", c.APIKey)
	}
}
