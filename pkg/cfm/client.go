// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cfm

import (
	"fmt"

	connectapi "github.com/HewlettPackard/hpecli/cfm/connectapi"
	"github.com/HewlettPackard/hpecli/internal/platform/rest"
)

const apiDefault = 800

// newOVClient creates a new OneView Client from username/password
// Creating our own constructor method that defaults
func newCFMClient(host, username, password string) *connectapi.CFMClient {
	return &connectapi.CFMClient{
		Host:     host,
		Username: username,
		Password: password,
		Token:    "",
	}
}

// newOVClientFromAPIKey creates a new OneView Client from existing API sessions key
func newCFMClientFromAPIKey(host, apikey string) *connectapi.CFMClient {
	return &connectapi.CFMClient{
		Host:     host,
		Token:    apikey,
		Username: "",
		Password: "",
	}
}

// login creates a OneView session
func login(host, username, password string) (string, error) {
	cfmClient, err := connectapi.GetAuthToken(host, username, password)
	if err == nil {
		return cfmClient.Token, nil
	}

	return "", fmt.Errorf("unable to session Token from login request")
}

// addAPIHeaders sets OneView API version
func addAPIHeaders() func(*rest.Request) {
	return func(r *rest.Request) {
		r.Header.Set("X-API-Version", "800")
	}
}
