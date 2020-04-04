// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cfm

import (
	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

const CFMAPIKeyPrefix = "hpecli_cfm_token_"
const CFMContextKey = "hpecli_cfm_context"

func hostAndToken() (host, token string, err error) {
	c := context.New(CFMContextKey)

	host, err = c.ModuleContext()
	if err != nil {
		return "", "", err
	}

	if err = c.HostData(dataKey(host), &token); err != nil {
		return "", "", err
	}

	return host, token, nil
}

func saveContextAndHostData(host, token string) error {
	c := context.New(CFMContextKey)
	if err := c.SetModuleContext(host); err != nil {
		return err
	}

	return c.SetHostData(dataKey(host), token)
}

func hostData(host string) (token string, err error) {
	c := context.New(CFMContextKey)
	if err = c.HostData(dataKey(host), &token); err != nil {
		return "", err
	}

	return token, nil
}

func getContext() (string, error) {
	c := context.New(CFMContextKey)
	return c.ModuleContext()
}

func setContext(host string) error {
	c := context.New(CFMContextKey)
	return c.SetModuleContext(host)
}

func deleteSavedHostData(host string) error {
	c := context.New(CFMContextKey)
	return c.DeleteHostData(dataKey(host))
}

func dataKey(apiEndpoint string) string {
	return CFMAPIKeyPrefix + apiEndpoint
}
