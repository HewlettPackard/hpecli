// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package update

import (
	"fmt"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/HewlettPackard/hpecli/pkg/version"
	"github.com/tcnksm/go-latest"
)

// json file that describes the latest release version.  Should be updated when new versions are published
// can alternatively change to using github tags once we real releases
const versionHost = "raw.github.com"
const versionPath = "/HewlettPackard/hpecli/master/site/release-version.json"

var versionURL = fmt.Sprintf("https://%s%s", versionHost, versionPath)

// cache retrieved response request
var response *latest.CheckResponse

// IsUpdateAvailable checks if a later version is avaialbe of the CLI binary
func IsUpdateAvailable() bool {
	// Since the CLI is a short-lived process.. only retrieve it once per instance.
	// If we haven't retrieved it, then go get it.  Once we get a copy, we then
	// store a cached copy and never get it again for the life of the CLI instance.
	if response != nil {
		return response.Outdated
	}
	json := &latest.JSON{
		URL: versionURL,
	}

	ver := version.Get()
	if ver == "" {
		ver = "0.0.0"
	}
	logger.Debug("Current version is: " + ver)
	logger.Debug("Checking for a newer version at: " + versionURL)
	res, err := checkUpdate(json, ver)
	if err != nil {
		logger.Debug("Unable to determine if a new version of the CLI is available")
		logger.Debug("Error: %v", err)
		return false
	}
	logger.Debug("Newest Version=" + res.Current)
	logger.Debug("Meta.Message=" + res.Meta.Message)
	logger.Debug("Meta.URL=" + res.Meta.URL)
	logger.Debug(fmt.Sprintf("outdate=%v", res.Outdated))
	logger.Debug(fmt.Sprintf("latest=%v", res.Latest))
	logger.Debug(fmt.Sprintf("new=%v", res.New))
	// cache a copy
	response = res
	return res.Outdated
}

func checkUpdate(s latest.Source, curVer string) (*latest.CheckResponse, error) {
	return latest.Check(s, curVer)
}
