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

func saveData(apiEndpoint, tenantID, token string) error {
	c := context.New(glContextKey)
	if err := c.SetModuleContext(apiEndpoint); err != nil {
		return err
	}

	key := dataKey(apiEndpoint)

	return c.SetHostData(key, &glContextData{key, token, tenantID})
}

func getData() (*glContextData, error) {
	c := context.New(glContextKey)

	apiEndpoint, err := c.ModuleContext()
	if err != nil {
		return nil, err
	}

	key := dataKey(apiEndpoint)

	var value glContextData
	if err = c.HostData(key, &value); err != nil {
		return nil, err
	}

	return &value, nil
}

func setContext(key string) error {
	c := context.New(glContextKey)
	return c.SetModuleContext(key)
}

func dataKey(apiEndpoint string) string {
	return glAPIKeyPrefix + apiEndpoint
}
