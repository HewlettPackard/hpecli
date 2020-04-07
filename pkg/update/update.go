// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package update

import (
	"fmt"
	"io"
	"net/http"

	goupdate "github.com/inconshreveable/go-update"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewUpdateCommand(version string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "update",
		Short: "Update hpe CLI",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runUpdate(version)
		},
	}

	return cmd
}

func runUpdate(localVersion string) error {
	resp, err := checkUpdate(&jsonSource{url: versionURL}, localVersion)
	if err != nil {
		return err
	}

	if !resp.UpdateAvailable {
		logrus.Warning("No update available.  No action taken")
		return nil
	}

	if err := downloadUpdate(resp); err != nil {
		return err
	}

	logrus.Infof("Successfully update the cli to version: %s", resp.RemoteVersion)

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
		logrus.Warningf("Unable to update to new version of the application: %v", err)
		return err
	}

	return nil
}

func getResponseBody(url string) (io.ReadCloser, error) {
	client := &http.Client{}
	// Create a new GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
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
