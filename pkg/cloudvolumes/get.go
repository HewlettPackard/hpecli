// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolumes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"io/ioutil"
	"crypto/tls"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/HewlettPackard/hpecli/pkg/store"

	"github.com/spf13/cobra"
)

var (
	path string
	jsonresults bool
)

func init() {
	cmdCloudVolumesGet.Flags().StringVar(&host, "host", "", "CloudVolumes portal")
	cmdCloudVolumesGet.Flags().StringVar(&path, "path", "", "resource to get")
	cmdCloudVolumesGet.MarkFlagRequired("path")
	cmdCloudVolumesGet.Flags().BoolVar(&jsonresults, "json", false, "display result in json")

}

// getCmd represents the get command
var cmdCloudVolumesGet= &cobra.Command{
	Use:   "get",
	Short: "Get from Cloud Volumes: hpecli cloudvolumes get",
	RunE:  runGet,
}

func runGet(cmd *cobra.Command, args []string) error {
	logger.Info("cloudvolumes/get called")

	if host == "" {
		return fmt.Errorf("must provide --host or -h")
	}
	//if !strings.HasPrefix(host, "http") {
	//	host = fmt.Sprintf("http://%s", host)
	//}
	db, err := store.Open()
	if err != nil {
		logger.Debug("unable to open keystore: %v", err)
		return fmt.Errorf("%v", err)
	}

	var token string
	if err := db.Get(key(), &token); err == store.ErrNotFound {
	// key not found
		logger.Debug("Key not found, not logged in")
		return fmt.Errorf("Key not found, not logged in")	
	} else if err != nil {
		logger.Debug("unable to open keystore: %v", err)
		return fmt.Errorf("%v", err)
	} 
	db.Close()

	logger.Info(fmt.Sprintf("Attempting login with user: %v, at: %v", username, host))

	switch path {
	case "volumes":
		url := "https://" + host + "/api/v2/cloud_volumes"
		request, _ := http.NewRequest("GET", url, nil)
		authstring := "username:" + token 
		encoded := base64.StdEncoding.EncodeToString([]byte(authstring))
		request.Header.Set("Authorization", "Basic " + encoded)

		// Ignore invalid certificate (should be done once globally)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		
		response, err := client.Do(request)
		if err != nil {
			logger.Debug("unable to get volumes %v", err)
			return fmt.Errorf("%v", err)
		}

		logger.Info(fmt.Sprintf("Status Code = %v", response.StatusCode))

		if response.StatusCode == 401 {
			logger.Debug("Stored token has expired")
			// Stored token has expired delete it from key store
			db, err := store.Open()
			if err != nil {
				logger.Debug("unable to open keystore: %v", err)
				return fmt.Errorf("%v", err)
			}
			db.Delete(key())
			db.Close()
			return fmt.Errorf("Stored token has expired, please login")	
		}
		if response.StatusCode == 200 {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				logger.Debug("unable to parse response %v", err)
				return fmt.Errorf("%v", err)
			}
			resstring := string(body)

			if jsonresults {
				// Just dump josn response
				fmt.Println(resstring)
			} else {
				// Let's parse the json and display some fields
				

			}
		}
	case "connections":

	default:
		logger.Debug("Unknown path %v", path)
		return fmt.Errorf("Unknown path %v", path)	
	}
return nil
}
