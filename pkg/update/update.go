// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package update

import (
	"fmt"
	"io"
	"net/http"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/HewlettPackard/hpecli/pkg/version"
	goupdate "github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"
)

// Cmd represents the update command
var Cmd = &cobra.Command{
	Use:   "update",
	Short: "Update the hpecli executable",
	RunE:  runUpdate,
}

func runUpdate(_ *cobra.Command, _ []string) error {
	localVer := version.Get()

	resp, err := checkUpdate(&jsonSource{url: versionURL}, localVer)
	if err != nil {
		return err
	}

	if !resp.UpdateAvailable {
		logger.Always("No update available.  No action taken")
		return nil
	}

	if err := downloadUpdate(resp); err != nil {
		return err
	}

	return nil
}

func downloadUpdate(cr *CheckResponse) error {
	body, err := getResponseBody(cr.URL)
	if err != nil {
		return err
	}
	defer body.Close()

	err = goupdate.Apply(body, goupdate.Options{Checksum: cr.CheckSum})
	if err != nil {
		logger.Warning("Unable to update to new version of the application: %v", err)
		return err
	}

	return nil
}

func getResponseBody(url string) (io.ReadCloser, error) {
	client := &http.Client{}
	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	//nolint:bodyclose // body is closed above in downloadUpdate
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("\"%s\" retrieving remote executable at: %v", resp.Status, url)
	}

	return resp.Body, nil
}
