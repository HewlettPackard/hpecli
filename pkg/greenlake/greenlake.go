// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"fmt"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/HewlettPackard/hpecli/pkg/store"
	"github.com/spf13/cobra"
)

const greenLakeAPIKeyPrefix = "hpecli_greenlake_token_"
const greenLakeTenantIDPrefix = "hpecli_greenlake_tenantid_"
const greenLakeContextKey = "hpecli_greenlake_context"

func init() {
	Cmd.AddCommand(glGetCmd)
	Cmd.AddCommand(glLoginCmd)
}

// Cmd represents the ilo command
var Cmd = &cobra.Command{
	Use:   "greenlake",
	Short: "Access to HPE GreenLake commands",
}

func getTokenTenantID() (host, tenantID, apiKey string) {
	db, err := store.Open()
	if err != nil {
		logger.Debug("unable to open keystore: %v", err)
		return "", "", ""
	}

	defer db.Close()

	var contextValue string
	if err := db.Get(greenLakeContextKey, &contextValue); err != nil {
		logger.Debug("Unable to retrieve current context.")
		return "", "", ""
	}

	apiKey = makeAPIKey(contextValue)

	var apiKeyValue string
	if err := db.Get(apiKey, &apiKeyValue); err != nil {
		logger.Debug("Unable to retrieve current context.")
		return "", "", ""
	}

	tenantID = makeTenantID(contextValue)

	var tenantIDValue string
	if err := db.Get(tenantID, &tenantIDValue); err != nil {
		logger.Debug("Unable to retrieve current context.")
		return "", "", ""
	}

	return contextValue, tenantIDValue, apiKeyValue
}

func setTokenTentanID(host, glTenantID, glAccessToken string) error {
	db, e := store.Open()
	if e != nil {
		return fmt.Errorf("unable to open keystore: %w", e)
	}

	defer db.Close()

	// Save context key
	if e := db.Put(greenLakeContextKey, host); e != nil {
		return fmt.Errorf("unable to retrieve current context because of %w", e)
	}

	// Save API Key
	apiKey := makeAPIKey(host)
	if e := db.Put(apiKey, glAccessToken); e != nil {
		return fmt.Errorf("unable to save apiKey for %s because of %w", host, e)
	}

	// Save TeantID Key
	tenantID := makeTenantID(host)
	if e := db.Put(tenantID, glTenantID); e != nil {
		return fmt.Errorf("unable to save tenantID for %s because of %w", host, e)
	}

	return nil
}

func makeAPIKey(host string) string {
	return fmt.Sprintf("%s%s", greenLakeAPIKeyPrefix, host)
}

func makeTenantID(host string) string {
	return fmt.Sprintf("%s%s", greenLakeTenantIDPrefix, host)
}
