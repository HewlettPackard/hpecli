// (C) Copyright 2019 Hewlett Packard Enterprise Development LP

package greenlake

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/rest"
)

// GLClient - wrapper class for greenlake api's
type GLClient struct {
	GrantType    string
	ClientID     string
	ClientSecret string
	TenantID     string
	Host         string
	APIKey       string
	*rest.Request
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

// newGLClient create
func newGLClient(grantType, clientID, secretKey, tenantID, host string) *GLClient {
	return &GLClient{
		GrantType:    grantType,
		ClientID:     clientID,
		ClientSecret: secretKey,
		TenantID:     tenantID,
		Host:         host,
		APIKey:       "",
	}
}

// newGLClientFromAPIKey creates a new GreenLake GLClient from existing API sessions key
func newGLClientFromAPIKey(host, tenantID, token string) *GLClient {
	return &GLClient{
		GrantType:    "client_credentials",
		ClientID:     "",
		ClientSecret: "",
		APIKey:       token,
		TenantID:     tenantID,
		Host:         host,
	}
}

// Login api
func (c *GLClient) login() (*sessionData, error) {
	const uriPath = "/identity/v1/token"

	sd := &sessionData{}

	loginJSON := fmt.Sprintf(`{"grant_type":"%s", "client_id":"%s",
	"client_secret":"%s", "tenant_id":"%s"}`,
		c.GrantType, c.ClientID, c.ClientSecret, c.TenantID)

	resp, err := rest.Post(c.Host+uriPath, strings.NewReader(loginJSON),
		rest.AddJSONMimeType(), rest.AllowSelfSignedCerts())
	if err != nil {
		return sd, err
	}

	if resp.StatusCode != http.StatusOK {
		return sd, fmt.Errorf("unable to create login sessions to Green Lake.  Response was: %+v", resp.Status)
	}

	var result Token

	err = resp.Unmarshall(&result)
	if err != nil {
		return sd, fmt.Errorf("nable to create login token from session")
	}

	if result.AccessToken == "" {
		return sd, fmt.Errorf("nable to create login token from session")
	}

	sd.Host = c.Host
	sd.Token = result.AccessToken
	sd.TenantID = c.TenantID

	return sd, nil
}

// users to list users
func (c *GLClient) users() ([]byte, error) {
	uriPath := fmt.Sprintf("/scim/v1/tenant/" + c.TenantID + "/" + "Users")

	resp, err := rest.Get(c.Host+uriPath, c.addAuth(c.APIKey), rest.AllowSelfSignedCerts())
	if err != nil {
		return []byte{}, err
	}

	return resp.JSON(), nil
}

// addAuth func
func (c *GLClient) addAuth(apiKey string) func(*rest.Request) {
	return func(r *rest.Request) {
		r.Header.Add("Authorization", "Bearer "+apiKey)
	}
}
