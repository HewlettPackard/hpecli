// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

const cvAPIKeyPrefix = "hpecli_cloudvolume_token_"
const cvContextKey = "hpecli_cloudvolume_context"

const cvDefaultHost = "https://demo.cloudvolumes.hpe.com"



func hostAndToken() (host, token string, err error) {
	c := context.New(cvContextKey)

	apiEndpoint, err := c.ModuleContext()
	if err != nil {
		return "", "", err
	}

	key := dataKey(apiEndpoint)
	if err = c.HostData(key, &token); err != nil {
		return "", "", err
	}

	return apiEndpoint, token, nil
}

func saveData(apiEndpoint, token string) error {
	c := context.New(cvContextKey)
	if err := c.SetModuleContext(apiEndpoint); err != nil {
		return err
	}

	key := dataKey(apiEndpoint)

	return c.SetHostData(key, token)
}

func dataKey(apiEndpoint string) string {
	return cvAPIKeyPrefix + apiEndpoint
}

func hostData(host string) (token string, err error) {
	c := context.New(cvContextKey)
	if err = c.HostData(dataKey(host), &token); err != nil {
		return "", err
	}

	return token, nil
}

func deleteSavedHostData(host string) error {
	c := context.New(cvContextKey)
	return c.DeleteHostData(dataKey(host))
}
