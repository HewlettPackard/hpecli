// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"github.com/HewlettPackard/hpecli/pkg/context"
)

const glAPIKeyPrefix = "hpecli_greenlake_token_"
const glContextKey = "hpecli_greenlake_context"

type glContextData struct {
	Host     string
	APIKey   string
	TenantID string
}

func storeContext(key, tenantID, token string) error {
	c := context.New(glContextKey, glAPIKeyPrefix)
	return c.SetAPIKey(key, &glContextData{key, token, tenantID})
}

func getContext() (*glContextData, error) {
	c := context.New(glContextKey, glAPIKeyPrefix)

	var d glContextData
	err := c.APIKey(&d)

	return &d, err
}

func changeContext(key string) error {
	c := context.New(glContextKey, glAPIKeyPrefix)
	return c.ChangeContext(key)
}
