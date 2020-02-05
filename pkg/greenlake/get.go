// (C) Copyright 2019 Hewlett Packard Enterprise Development LP

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

const users = "users"

func init() {
	glGetCmd.Flags().StringVar(&getPath, "path", "p", "path to a RedFish item")
	glGetCmd.Flags().BoolVar(&getJSONResult, "json", false, "display result in json")
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

	host, tenantID, apiKey := getTokenTenantID()
	glc := NewGLClientFromAPIKey(host, tenantID, apiKey)

	switch getPath {
	case users:
		body, err := glc.GetUsers("Users")

		if err != nil {
			logger.Debug("unable to get the users with the supplied credentials: %v", err)
			return err
		}

		if getJSONResult {
			resstring := string(body)
			fmt.Println(resstring)
		} else {
			var result []User
			if err := json.Unmarshal(body, &result); err != nil {
				return err
			}
			for index := range result {
				fmt.Printf("Name: %s : Email: %s Active: %t\n", result[index].DisplayName,
					result[index].UserName, result[index].Active)
			}
		}

	default:
		fmt.Println("Unknown path: ", getPath)
	}

	return nil
}
