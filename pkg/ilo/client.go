// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/logger"
)

type ILOClient struct {
	restClient
}

func NewILOClient(host, username, password string) *ILOClient {
	return &ILOClient{
		restClient{
			Endpoint: host,
			Username: username,
			Password: password,
		},
	}
}

func NewILOClientFromAPIKey(host, token string) *ILOClient {
	return &ILOClient{
		restClient{
			Endpoint: host,
			APIKey:   token,
		},
	}
}

func (c *ILOClient) Login() (string, error) {
	const uriPath = "/redfish/v1/SessionService/Sessions/"

	loginJSON := fmt.Sprintf(`{"UserName":"%s", "Password":"%s"}`, c.Username, c.Password)

	data, err := c.restAPICall("POST", uriPath, strings.NewReader(loginJSON))
	if err != nil {
		logger.Critical("Unable to login as: %s to: %s", c.Username, c.Endpoint)
		return "", err
	}

	return string(data), nil
}

func (c *ILOClient) GetServiceRoot() ([]byte, error) {
	const uriPath = "/redfish/v1/"

	data, err := c.restAPICall("GET", uriPath, nil)
	if err != nil {
		logger.Critical("Unable to get service root at: %s", c.Endpoint)
		return nil, err
	}

	return data, nil
}
