// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"github.com/HewlettPackard/hpecli/pkg/context"
	"github.com/HewlettPackard/hpecli/pkg/db"
)

const greenLakeAPIKeyPrefix = "hpecli_greenlake_token_"
const greenLakeTenantIDPrefix = "hpecli_greenlake_tenantid_"
const greenLakeContextKey = "hpecli_greenlake_context"

type glContextData struct {
	host     string
	tenantID string
}

func storeContext(key, token string) error {
	c := context.New(cvContextKey, cvAPIKeyPrefix, db.Open)
	return c.SetAPIKey(key, &glContextData{key, token})
}

func getContext() (*glContextData, error) {
	c := context.New(cvContextKey, cvAPIKeyPrefix, db.Open)

	var d glContextData
	err := c.APIKey(&d)

	return &d, err
}
