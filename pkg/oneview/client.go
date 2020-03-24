// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/rest"
	"github.com/HewlettPackard/oneview-golang/ov"
	ovrest "github.com/HewlettPackard/oneview-golang/rest"
)

const apiDefault = 800

// newOVClient creates a new OneView Client from username/password
// Creating our own constructor method that defaults
func newOVClient(host, username, password string) *ov.OVClient {
	return &ov.OVClient{
		Client: ovrest.Client{
			User:       username,
			Password:   password,
			Domain:     "LOCAL",
			APIKey:     "",
			APIVersion: apiDefault,
			SSLVerify:  true,
			Endpoint:   host,
			IfMatch:    "",
		},
	}
}

// newOVClientFromAPIKey creates a new OneView Client from existing API sessions key
func newOVClientFromAPIKey(host, apikey string) *ov.OVClient {
	return &ov.OVClient{
		Client: ovrest.Client{
			User:       "",
			Password:   "",
			Domain:     "LOCAL",
			APIKey:     apikey,
			APIVersion: apiDefault,
			SSLVerify:  true,
			Endpoint:   host,
			IfMatch:    "",
		},
	}
}

// login creates a OneView session
func login(host, username, password string) (string, error) {
	const uriPath = "/rest/login-sessions"

	loginJSON := fmt.Sprintf(`{"userName":"%s", "password":"%s", "authLoginDomain":"LOCAL", "loginMsgAck":"true"}`,
		username, password)

	opts := func(r *rest.Request) {
		rest.AddJSONMimeType()(r)
		rest.AllowSelfSignedCerts()(r)
		addAPIHeaders()(r)
	}

	resp, err := rest.Post(host+uriPath, strings.NewReader(loginJSON), opts)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to create login sessions to oneview.  Response was: %+v", resp.Status)
	}

	var session ov.Session

	err = resp.Unmarshall(&session)
	if err == nil {
		return session.ID, nil
	}

	return "", fmt.Errorf("unable to session Token from login request")
}

// addAPIHeaders sets OneView API version
func addAPIHeaders() func(*rest.Request) {
	return func(r *rest.Request) {
		r.Header.Set("X-API-Version", "800")
	}
}
