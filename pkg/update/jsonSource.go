// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package update

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/hashicorp/go-version"
)

type jsonSource struct {
	url string
}

type jsonResponse struct {
	Version   string `json:"version"`
	Message   string `json:"message"`
	UpdateURL string `json:"url"`
	PublicKey string `json:"publickey"`
	CheckSum  string `json:"checksum"`
}

func (j *jsonSource) validate() error {
	if j.url == "" {
		return fmt.Errorf("remote URL must be set")
	}

	if _, err := url.Parse(j.url); err != nil {
		return fmt.Errorf("%s is invalid URL: %s", j.url, err.Error())
	}

	return nil
}

func (j *jsonSource) get() (*remoteResponse, error) {
	client := &http.Client{}
	// Create a new GET request
	req, err := http.NewRequest(http.MethodGet, j.url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	jsonResponse, err := decodeResponse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to decode json response: %v", err)
	}

	result, err := mapResult(jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("unable to map json to remoteResponse: %v", err)
	}

	return result, nil
}

func decodeResponse(r io.Reader) (*jsonResponse, error) {
	result := &jsonResponse{}
	dec := json.NewDecoder(r)

	if err := dec.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func mapResult(j *jsonResponse) (*remoteResponse, error) {
	v, err := version.NewVersion(j.Version)
	if err != nil {
		return nil, fmt.Errorf("did not get a remote version. %v", err)
	}

	downloadedPublicKey := decodeField("PublicKey", j.PublicKey)
	downloadedCheckSum := decodeField("CheckSum", j.CheckSum)

	return &remoteResponse{
		version:   v,
		message:   j.Message,
		updateURL: j.UpdateURL,
		publicKey: downloadedPublicKey,
		checkSum:  downloadedCheckSum,
	}, nil
}

// we need a nil value rather than an empty byte array.  the file download verification
// treats non-nil values as values to compare against.
func decodeField(name, value string) []byte {
	if value == "" {
		return nil
	}

	result, err := hex.DecodeString(value)
	if err != nil {
		logger.Warning("Unable to decode remote %s field. "+
			"It will be ignored and not used to verify the remote content", name)
		logger.Debug("Problem Field: %s=%v", name, value)

		return nil
	}

	return result
}
