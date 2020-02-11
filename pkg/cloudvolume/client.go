// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/logger"
)

type CVClient struct {
	restClient
}

func NewCVClient(host, username, password string) *CVClient {
	return &CVClient{
		restClient{
			Endpoint: host,
			Username: username,
			Password: password,
		},
	}
}

func (c CVClient) Login() (string, error) {
	const uriPath = "/auth/login"

	loginJSON := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, c.Username, c.Password)

	data, err := c.restAPICall("POST", uriPath, strings.NewReader(loginJSON))
	if err != nil {
		logger.Critical("Unable to login as: %s to: %s", c.Username, c.Endpoint)
		return "", err
	}

	type loginResponse struct {
		Geo   string `json:"geo"`
		Token string `json:"token"`
	}

	var lr loginResponse

	if err := json.Unmarshal(data, &lr); err != nil {
		logger.Debug("expcted login response, but received: %s", lr)
		return "", err
	}

	return lr.Token, nil
}
