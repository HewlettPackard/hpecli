//(C) Copyright 2019 Hewlett Packard Enterprise Development LP

package greenlake

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/greenlake/client"
	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/HewlettPackard/hpecli/pkg/store"
	"github.com/spf13/cobra"
)

// User structure
type User struct {
	Active      bool   `json:"active"`
	DisplayName string `json:"displayName"`
	UserName    string `json:"userName"`
	Name        Name   `json:"name"`
}

// Name structure
type Name struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
}

var (
	getPath       string
	getJSONResult bool
)

func init() {
	cmdGreenLakeGet.Flags().StringVar(&getPath, "path", "p", "path to a RedFish item")
	cmdGreenLakeGet.MarkFlagRequired("path")
	cmdGreenLakeGet.Flags().StringVar(&host, "host", "", "greenlake ip address")
	cmdGreenLakeGet.Flags().BoolVar(&getJSONResult, "json", false, "display result in json")

}

// glhc/getCmd represents the glhc/get command
var cmdGreenLakeGet = &cobra.Command{
	Use:   "get",
	Short: "A greenlake get command description",
	RunE:  runGet,
}

func runGet(cmd *cobra.Command, args []string) error {
	logger.Info("greenlake/get called")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	if !strings.HasPrefix(host, "http") {
		host = fmt.Sprintf("http://%s", host)
	}

	switch getPath {
	case "users":
		var (
			token, tenantID string
		)
		db, err := store.Open()
		if err != nil {
			logger.Debug("unable to open keystore: %v", err)
			return fmt.Errorf("%v", err)
		}
		key1 := "hpecli_greenlake_token_" + host
		key2 := "hpecli_greenlake_tenantid_" + host
		db.Get(key1, &token)
		db.Get(key2, &tenantID)
		db.Close()

		//fmt.Printf("key1=%s Token=%s\n", key1, token)
		//fmt.Printf("key2=%s  TenantId=%s\n", key2, tenantID)

		if token == "" {
			fmt.Println("Not logged in!")
		} else {
			client := client.NewGreenLakeClient("", "", "", tenantID, host)
			body, err := client.GetUsers("Users", token)
			if err != nil {
				logger.Debug("unable to get the token with the supplied credentials: %v", err)
				return err
			}
			var result []User
			json.Unmarshal(body, &result)
			resstring := string(body)
			if getJSONResult {
				fmt.Println(resstring)
			} else {
				for _, user := range result {
					fmt.Printf("Name: %s : Email: %s Active: %t\n", user.DisplayName, user.UserName, user.Active)
				}
			}
		}
	default:
		fmt.Println("Unknown path: ", getPath)
	}
	return nil
}
