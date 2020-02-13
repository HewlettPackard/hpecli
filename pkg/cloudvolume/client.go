// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/internal/rest"
)

type CVClient struct {
	Host     string
	Username string
	Password string
	APIKey   string
	*rest.Request
}

func NewCVClient(host, username, password string) *CVClient {
	return &CVClient{
		Host:     host,
		Username: username,
		Password: password,
	}
}

func NewCVClientFromAPIKey(host, token string) *CVClient {
	return &CVClient{
		Host:   host,
		APIKey: token,
	}
}

func (c *CVClient) Login() (string, error) {
	const uriPath = "/auth/login"

	postBody := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, c.Username, c.Password)

	resp, err := rest.Post(c.Host+uriPath, strings.NewReader(postBody), AddJSONMimeType())
	if err != nil {
		return "", err
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

func (c *CVClient) GetCloudVolumes() ([]byte, error) {
	const uriPath = "/api/v2/cloud_volumes"

	resp, err := rest.Get(c.Host+uriPath, c.AddAuth())
	if err != nil {
		return []byte{}, err
	}

	return resp.JSON(), nil
}

func (c *CVClient) AddAuth() func(*rest.Request) {
	return func(r *rest.Request) {
		r.Request.SetBasicAuth("username", c.APIKey)
	}
}

func AddJSONMimeType() func(*rest.Request) {
	return func(r *rest.Request) {
		r.Header.Set("Content-Type", "application/json")
	}
}
