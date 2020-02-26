// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"github.com/HewlettPackard/hpecli/pkg/context"
)

const ovAPIKeyPrefix = "hpecli_oneview_token_"
const ovContextKey = "hpecli_oneview_context"

func hostAndToken() (host, token string, err error) {
	c := context.New(ovContextKey)

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
	c := context.New(ovContextKey)
	if err := c.SetModuleContext(host); err != nil {
		return err
	}

	return c.SetHostData(dataKey(host), token)
}

func hostData(host string) (token string, err error) {
	c := context.New(ovContextKey)
	if err = c.HostData(dataKey(host), &token); err != nil {
		return "", err
	}

	return token, nil
}

func setContext(host string) error {
	c := context.New(ovContextKey)
	return c.SetModuleContext(host)
}

func deleteSavedHostData(host string) error {
	c := context.New(ovContextKey)
	return c.DeleteHostData(dataKey(host))
}

func dataKey(apiEndpoint string) string {
	return ovAPIKeyPrefix + apiEndpoint
}
