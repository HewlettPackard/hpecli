// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/internal/rest"
	"github.com/HewlettPackard/hpecli/pkg/logger"
)

type Client struct {
	Host     string
	Username string
	Password string
	APIKey   string
	*rest.Request
}

func NewILOClient(host, username, password string) *Client {
	return &Client{
		Host:     host,
		Username: username,
		Password: password,
	}
}

func NewILOClientFromAPIKey(host, token string) *Client {
	return &Client{
		Host:   host,
		APIKey: token,
	}
}

func (c *Client) Login() (string, string, error) {
	const uriPath = "/redfish/v1/sessionservice/sessions/"

	loginJSON := fmt.Sprintf(`{"UserName":"%s", "Password":"%s"}`, c.Username, c.Password)

	resp, err := rest.Post(c.Host+uriPath, strings.NewReader(loginJSON),
		rest.AddJSONMimeType(), rest.AllowSelfSignedCerts())
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("unable to create login sessions to ilo.  Repsponse was: %+v", resp.Status)
	}

	token := resp.Header.Get("X-Auth-Token")
	location := resp.Header.Get("Location")

	if token == "" {
		return "", "", fmt.Errorf("unable to create login toekn from session")
	}

	return token, location, nil
}

func (c *Client) Logout(sessionLocation string) error {
	resp, err := rest.Delete(sessionLocation, AddAuth(c.APIKey), rest.AllowSelfSignedCerts())
	if err != nil {
		return err
	}

	// 2xx status codes are OK
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logger.Debug("logout failed:  %+v", resp.Status)
		return fmt.Errorf("unable to successfully logout of iLO: %s", c.Host)
	}
	return nil
}

func (c *Client) GetServiceRoot() ([]byte, error) {
	const uriPath = "/redfish/v1/"

	resp, err := rest.Get(c.Host+uriPath, AddAuth(c.APIKey), rest.AllowSelfSignedCerts())
	if err != nil {
		return []byte{}, err
	}

	return resp.JSON(), nil
}

func AddAuth(apiKey string) func(*rest.Request) {
	return func(r *rest.Request) {
		r.Header.Add("X-Auth-Token", apiKey)
	}
}
