// (C) Copyright 2019 Hewlett Packard Enterprise Development LP

package greenlake

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/internal/rest"
)

// GLClient - wrapper class for greenlake api's
type GLClient struct {
	GrantType    string
	ClientID     string
	ClientSecret string
	TenantID     string
	Host         string
	APIKey       string
}

// Token structure
type Token struct {
	AccessToken     string `json:"access_token"`
	Scope           string `json:"scope"`
	TokenType       string `json:"token_type"`
	Expiry          string `json:"expiry"`
	ExpiresIn       int    `json:"expires_in"`
	AccessTokenOnly bool   `json:"accessTokenOnly"`
}

// User structure
type User struct {
	Active      bool   `json:"active"`
	DisplayName string `json:"displayName"`
	UserName    string `json:"userName"`
	Name        Name   `json:"name"`
}

// Name structure
type Name struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
}

// NewGLClient create
func NewGLClient(grantType, clientID, secretKey, tenantID, host string) *GLClient {
	return &GLClient{
		GrantType:    grantType,
		ClientID:     clientID,
		ClientSecret: secretKey,
		TenantID:     tenantID,
		Host:         host,
		APIKey:       "",
	}
}

// NewGLClientFromAPIKey creates a new GreenLake GLClient from existing API sessions key
func NewGLClientFromAPIKey(host, tenantID, apikey string) *GLClient {
	return &GLClient{
		GrantType:    "client_credentials",
		ClientID:     "",
		ClientSecret: "LOCAL",
		APIKey:       apikey,
		TenantID:     tenantID,
		Host:         host,
	}
}

// GetToken api
func (c *GLClient) GetToken() (Token, error) {
	const uriPath = "/identity/v1/token"

	postBody := fmt.Sprintf(`{"grant_type":"%s", "client_id":"%s", "client_secret":"%s", "tenant_id":"%s"}`, c.GrantType, c.Password, c.ClientSecret, c.TenantID)

	resp, err := rest.Post(c.Host+uriPath, strings.NewReader(postBody), AddJSONMimeType())
	if err != nil {
		return "", err
	}

	var result Token

	err = resp.Unmarshall(&result)
	if err == nil {
		return result.AccessToken, nil
	}

	return "", fmt.Errorf("unable to get response from login command")
}

// GetUsers to list users
func (c *GLClient) GetUsers(path string) ([]byte, error) {
	uriPath := fmt.Sprintf("/scim/v1/tenant/" + c.TenantID + "/" + path)

	resp, err := rest.Get(c.Host+uriPath, c.AddAuth())
	if err != nil {
		return []byte{}, err
	}

	return resp.JSON(), nil
}

// AddAuth func
func (c *GLClient) AddAuth() func(*rest.Request) {
	return func(r *rest.Request) {
		r.Request.SetBasicAuth("username", c.APIKey)
	}
}

// AddJSONMimeType func
func AddJSONMimeType() func(*rest.Request) {
	return func(r *rest.Request) {
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept", "pplication/scim+json")
	}
}
