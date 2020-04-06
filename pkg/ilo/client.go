// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/rest"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Host     string
	Username string
	Password string
	APIKey   string
	*rest.Request
}

func newILOClient(host, username, password string) *Client {
	return &Client{
		Host:     host,
		Username: username,
		Password: password,
	}
}

func newILOClientFromAPIKey(host, token string) *Client {
	return &Client{
		Host:   host,
		APIKey: token,
	}
}

func (c *Client) login() (*sessionData, error) {
	const uriPath = "/redfish/v1/sessionservice/sessions/"

	sd := &sessionData{}

	loginJSON := fmt.Sprintf(`{"UserName":"%s", "Password":"%s"}`, c.Username, c.Password)

	resp, err := rest.Post(c.Host+uriPath, strings.NewReader(loginJSON),
		rest.AddJSONMimeType(), rest.AllowSelfSignedCerts())
	if err != nil {
		return sd, err
	}

	if resp.StatusCode != http.StatusCreated {
		return sd, fmt.Errorf("unable to create login sessions to iLO.  Response was: %+v", resp.Status)
	}

	token := resp.Header.Get("X-Auth-Token")
	location := resp.Header.Get("Location")

	if token == "" {
		return sd, fmt.Errorf("unable to create login toekn from session")
	}

	sd.Host = c.Host
	sd.Token = token
	sd.Location = location

	return sd, nil
}

func (c *Client) logout(sessionLocation string) error {
	resp, err := rest.Delete(sessionLocation, AddAuth(c.APIKey), rest.AllowSelfSignedCerts())
	if err != nil {
		return err
	}

	// 2xx status codes are OK
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logrus.Debugf("logout failed: %+v", resp.Status)
		return fmt.Errorf("unable to successfully logout of iLO: %s", c.Host)
	}

	return nil
}

func (c *Client) getServiceRoot() ([]byte, error) {
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
