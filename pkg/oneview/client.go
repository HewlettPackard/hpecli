// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"github.com/HewlettPackard/oneview-golang/ov"
	"github.com/HewlettPackard/oneview-golang/rest"
)

// NewOVClient creates a new OneView Client from username/password
// Creating our own constructor method that defaults
func NewOVClient(host string, username string, password string) *ov.OVClient {
	return &ov.OVClient{
		Client: rest.Client{
			//Method:     rest.GET,
			User:       username,
			Password:   password,
			Domain:     "LOCAL",
			APIKey:     "",
			APIVersion: 800,
			SSLVerify:  true,
			Endpoint:   host,
			IfMatch:    "",
			//Option:     rest.Options{},
		},
	}
}

// NewOVClientFromAPIKey creates a new OneView Client from existing API sessions key
func NewOVClientFromAPIKey(host, apikey string) *ov.OVClient {
	return &ov.OVClient{
		Client: rest.Client{
			//Method:     rest.GET,
			User:       "",
			Password:   "",
			Domain:     "LOCAL",
			APIKey:     apikey,
			APIVersion: 800,
			SSLVerify:  true,
			Endpoint:   host,
			IfMatch:    "",
			//Option:     rest.Options{},
		},
	}
}
