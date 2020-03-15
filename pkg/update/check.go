// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package update

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	gover "github.com/hashicorp/go-version"
	"github.com/sirupsen/logrus"
)

// CheckResponse is a response for a Check request.
type CheckResponse struct {
	// Remote repo has a newer version than the running version
	UpdateAvailable bool

	// Latest version available in remote repository
	RemoteVersion string

	// Messagse about the latest updated
	Message string

	// URL where the update can be downloaded
	URL string

	// If the download has been signed, this is the public key
	// that can be used to verify the signature
	PublicKey []byte

	// SHA265 Has of the download that can be used to validate
	// the integrity of the file after it was downloaded
	CheckSum []byte
}

type remoteResponse struct {
	version   *gover.Version
	message   string
	updateURL string
	publicKey []byte
	checkSum  []byte
}

type source interface {
	validate() error
	get() (*remoteResponse, error)
}

// EnvDisableUpdateCheck is an environmental variable to disable remote
// http request to check if a newer version of the CLI is available
const EnvDisableUpdateCheck = "HPECLI_DISABLE_UPDATE_CHECK"

// json file that describes the latest release version.  Should be updated when new versions are published
// can alternatively change to using github tags once we real releases
const versionHost = "raw.githubusercontent.com"

const versionPath = "/HewlettPackard/hpecli/Didier/MultiOSSupport/site/published-version.json"


var versionURL = fmt.Sprintf("https://%s%s", versionHost, versionPath)

var cacheResponse *CheckResponse

// CheckForUpdate returns data about the availability of an updated version of the CLI
func CheckForUpdate(localVersion string) (*CheckResponse, error) {
	logrus.Debug("Local version is: " + localVersion)
	logrus.Debug("Checking for a newer version at: " + versionURL)

	res, err := checkUpdate(&jsonSource{url: versionURL}, localVersion)
	if err != nil {
		logrus.Debug("Unable to determine if a new version of the CLI is available")
		logrus.Debugf("Error: %v", err)

		return &CheckResponse{}, err
	}

	logrus.Debugf("%#v", res)

	return res, nil
}

// Check fetches last version information from its source
// and compares with target and return result (CheckResponse).
func checkUpdate(s source, lVersion string) (*CheckResponse, error) {
	// don't check if env var is setup to skip
	if os.Getenv(EnvDisableUpdateCheck) != "" {
		logrus.Debugf("%s set.  Not performing remote check", EnvDisableUpdateCheck)
		return &CheckResponse{}, nil
	}

	// Since the CLI is a short-lived process cache a copy and
	// return the cached copy if we have already retrieved the
	// results this session.
	if cacheResponse != nil {
		logrus.Debug("cacheResponse present.  Not making additional remote check")
		return cacheResponse, nil
	}

	localVersion, err := gover.NewVersion(lVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s, %s", lVersion, err.Error())
	}

	if err = s.validate(); err != nil {
		return nil, err
	}

	resp, err := s.get()
	if err != nil {
		return nil, err
	}

	var updateAvailable bool
	// If target > current, then update is available
	if resp.version.GreaterThan(localVersion) {
		updateAvailable = true
	}

	// retrieve GOOS from runtime GO variables
	osenv := runtime.GOOS

	if osenv == "" {
		logrus.Debugf("Runtime variable GOOS not set.  Not performing remote check")
		return &CheckResponse{}, nil
	} 

	// Substitute {{$GOOS}} with osenv in updateURL 
	resp.updateURL = strings.Replace(resp.updateURL, "{{$GOOS}}", osenv, 1)
	
	// windows uses a .exe filename while linux and MacOs don't
	var dotexe string = ""
	if osenv == "windows" {
		dotexe = ".exe" 			
	}

	logrus.Debugf("$GOOS=%v, $EXE=%v", osenv, dotexe)

	resp.updateURL = strings.Replace(resp.updateURL, "{{$EXE}}", dotexe, 1) 

	cacheResponse = &CheckResponse{
		UpdateAvailable: updateAvailable,
		RemoteVersion:   resp.version.String(),
		Message:         resp.message,
		URL:             resp.updateURL,
		PublicKey:       resp.publicKey,
		CheckSum:        resp.checkSum,
	}

	return cacheResponse, nil
}
