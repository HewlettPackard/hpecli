//(C) Copyright 2019 Hewlett Packard Enterprise Development LP

package greenlake

import (
	"encoding/json"
	"fmt"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	getPath       string
	getJSONResult bool
)

func init() {
	glGetCmd.Flags().StringVar(&getPath, "path", "p", "path to a RedFish item")
	// glGetCmd.Flags().StringVar(&glLoginData.host, "host", "", "greenlake ip address")
	glGetCmd.Flags().BoolVar(&getJSONResult, "json", false, "display result in json")
	// _ = glGetCmd.MarkFlagRequired("host")
	_ = glGetCmd.MarkFlagRequired("path")

}

// glhc/getCmd represents the glhc/get command
var glGetCmd = &cobra.Command{
	Use:   "get",
	Short: "A greenlake get command description",
	RunE:  runGlGet,
}

func runGlGet(_ *cobra.Command, _ []string) error {
	logger.Info("greenlake/get called")
	//cmd.SilenceUsage = true
	//cmd.SilenceErrors = true
	// if !strings.HasPrefix(glLoginData.host, "http") {
	// 	glLoginData.host = fmt.Sprintf("http://%s", glLoginData.host)
	// }
	host, tenantID, apiKey := getTokenTenantID()
	println("host %s api key is : ", host, apiKey)
	glc := NewGLClientFromAPIKey(host, tenantID, apiKey)
	println("NewGLClientFromAPIKey key is host %s TenantID %s accesstoken %s : ", host, tenantID, apiKey)

	switch getPath {
	case "users":
		body, err := glc.GetUsers("Users")
		println("body  is : ", body)
		if err != nil {
			logger.Debug("unable to get the users with the supplied credentials: %v", err)
			return err
		}
		var result []User
		if err := json.Unmarshal(body, &result); err != nil {
			return err
		}
		resstring := string(body)
		println("resstring  is : ", resstring)
		if getJSONResult {
			fmt.Println(resstring)
		} else {
			for _, user := range result {
				fmt.Printf("Name: %s : Email: %s Active: %t\n", user.DisplayName, user.UserName, user.Active)
			}
		}

	default:
		fmt.Println("Unknown path: ", getPath)
	}
	return nil
}
