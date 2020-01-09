// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolumes

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/HewlettPackard/hpecli/pkg/store"
	"github.com/spf13/cobra"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"bytes"
	"crypto/tls"
)

var (
	host     string
	username string
	password string
	token	 string
)

func init() {
	cmdCloudVolumesLogin.Flags().StringVar(&host, "host", "", "ilo ip address")
	cmdCloudVolumesLogin.Flags().StringVarP(&username, "username", "u", "", "ilo username")
	cmdCloudVolumesLogin.Flags().StringVarP(&password, "password", "p", "", "ilo passowrd")
}

// getCmd represents the get command
var cmdCloudVolumesLogin = &cobra.Command{
	Use:   "login",
	Short: "Login to Cloud Volumes: hpecli ilo login",
	RunE:  runLogin,
}

func runLogin(cmd *cobra.Command, args []string) error {
	fmt.Println("cloudvolumes/login called")

	if host == "" {
		return fmt.Errorf("must provide --host or -h")
	}
	//if !strings.HasPrefix(host, "http") {
	//	host = fmt.Sprintf("http://%s", host)
	//}

	if username == "" {
		// return fmt.Errorf("must provide --username or -u")
		fmt.Print("username: ")
		//var input string
		fmt.Scanln(&username)
		//fmt.Print(input)
	}

	if password == "" {
		// this really isn't secure to provide on the command line
		// need to provide a way to read from stdin
		//return fmt.Errorf("must provide --password or -p")
		fmt.Print("password: ")
		pass, _ := gopass.GetPasswd()
		password = string(pass)
	}
	
	fmt.Println("username", username)
	fmt.Println("password", password)

	fmt.Println(fmt.Sprintf("Attempting login with user: %v, at: %v", username, host))

	url := "https://" + host + "/auth/login"

	println("url:", url)
	jsonData := map[string]string{"email": username, 
								  "password": password}
	jsonValue, _ := json.Marshal(jsonData)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")

	// Ignore invalid certificate
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	// println("ready to call")
	response, err := client.Do(request)
	println("called")
	println("response", response.StatusCode)
	if err != nil {
		logger.Debug("unable to login %v", err)
		return fmt.Errorf("%v", err)
	}

	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logger.Debug("unable to login %v", err)
			return fmt.Errorf("%v", err)
		}
		var result map[string]string

		json.Unmarshal(body, &result)
		token = string(result["token"])

		// fmt.Println(token, location)
		fmt.Println("Logged in", token)
	}

	
	db, err := store.Open()
	if err != nil {
		logger.Debug("unable to open keystore: %v", err)
		return fmt.Errorf("%v", err)
	}
	defer db.Close()
	if err := db.Put(key(), &token); err != nil {
		return fmt.Errorf("%v", err)
	}
	db.Close()

	// if err != nil {
	// 	logger.Debug("Unable to login must provide --password or -p because of: %v", err)
	// 	return fmt.Errorf("Unable to loginmust provide --password or -p")
	// }

	return nil

}

func key() string {
	return fmt.Sprintf("hpecli_cloudvolumes_token_%s", host)
}
