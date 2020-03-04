// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

const glAPIKeyPrefix = "hpecli_greenlake_token_"
const glContextKey = "hpecli_greenlake_context"

type sessionData struct {
	Host     string
	Token    string
	TenantID string
}

func defaultSessionData() (data *sessionData, err error) {
	data = &sessionData{}
	c := context.New(glContextKey)

	host, err := c.ModuleContext()
	if err != nil {
		return data, err
	}

	if err = c.HostData(dataKey(host), &data); err != nil {
		return data, err
	}

	return data, nil
}

func saveContextAndSessionData(data *sessionData) error {
	c := context.New(glContextKey)
	if err := c.SetModuleContext(data.Host); err != nil {
		return err
	}

	return c.SetHostData(dataKey(data.Host), data)
}

func setContext(host string) error {
	c := context.New(glContextKey)
	return c.SetModuleContext(host)
}

func dataKey(apiEndpoint string) string {
	return glAPIKeyPrefix + apiEndpoint
}
