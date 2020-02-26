// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"github.com/HewlettPackard/hpecli/pkg/context"
)

const iloAPIKeyPrefix = "hpecli_ilo_token_"
const iloContextKey = "hpecli_ilo_context"

func hostAndToken() (host, token string, err error) {
	c := context.New(iloContextKey)

	host, err = c.ModuleContext()
	if err != nil {
		return "", "", err
	}

	if err = c.HostData(dataKey(host), &token); err != nil {
		return "", "", err
	}

	return host, token, nil
}

func saveData(host, token string) error {
	c := context.New(iloContextKey)
	if err := c.SetModuleContext(host); err != nil {
		return err
	}

	return c.SetHostData(dataKey(host), token)
}

func setContext(host string) error {
	c := context.New(iloContextKey)
	return c.SetModuleContext(host)
}

func deleteSavedData(host string) error {
	c := context.New(iloContextKey)
	return c.DeleteHostData(dataKey(host))
}

func dataKey(apiEndpoint string) string {
	return iloAPIKeyPrefix + apiEndpoint
}
