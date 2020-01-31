// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"fmt"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/HewlettPackard/hpecli/pkg/store"
	"github.com/spf13/cobra"
)

const oneViewAPIKeyPrefix = "hpecli_oneview_token_"
const oneViewContextKey = "hpecli_oneview_context"

func init() {
	Cmd.AddCommand(ovGetCmd)
	Cmd.AddCommand(ovLoginCmd)
}

// Cmd represents the ilo command
var Cmd = &cobra.Command{
	Use:   "oneview",
	Short: "Access to HPE OneView commands",
}

func apiKey() (host, apiKey string) {
	db, err := store.Open()
	if err != nil {
		logger.Debug("unable to open keystore: %v", err)
		return "", ""
	}

	defer db.Close()

	var contextValue string
	if err := db.Get(oneViewContextKey, &contextValue); err != nil {
		logger.Debug("Unable to retrieve current context.")
		return "", ""
	}

	apiKey = makeAPIKey(contextValue)

	var apiKeyValue string
	if err := db.Get(apiKey, &apiKeyValue); err != nil {
		logger.Debug("Unable to retrieve current context.")
		return "", ""
	}

	return contextValue, apiKeyValue
}

func setAPIKey(host, ovSessionID string) error {
	db, e := store.Open()
	if e != nil {
		return fmt.Errorf("unable to open keystore: %w", e)
	}

	defer db.Close()

	// Save context key
	if e := db.Put(oneViewContextKey, host); e != nil {
		return fmt.Errorf("unable to retrieve current context because of %w", e)
	}

	// Save API Key
	apiKey := makeAPIKey(host)
	if e := db.Put(apiKey, ovSessionID); e != nil {
		return fmt.Errorf("unable to save apiKey for %s because of %w", host, e)
	}

	return nil
}

func makeAPIKey(host string) string {
	return fmt.Sprintf("%s%s", oneViewAPIKeyPrefix, host)
}
