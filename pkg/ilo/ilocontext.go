// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

const iloAPIKeyPrefix = "hpecli_ilo_token_"
const iloContextKey = "hpecli_ilo_context"

type sessionData struct {
	Host     string
	Token    string
	Location string
}

func defaultSessionData() (data *sessionData, err error) {
	data = &sessionData{}
	c := context.New(iloContextKey)

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
	c := context.New(iloContextKey)
	if err := c.SetModuleContext(data.Host); err != nil {
		return err
	}

	return c.SetHostData(dataKey(data.Host), data)
}

func getSessionData(host string) (data *sessionData, err error) {
	data = &sessionData{}
	c := context.New(iloContextKey)

	if err = c.HostData(dataKey(host), &data); err != nil {
		return data, err
	}

	return data, nil
}

func setContext(host string) error {
	c := context.New(iloContextKey)
	return c.SetModuleContext(host)
}

func deleteSessionData(host string) error {
	c := context.New(iloContextKey)
	return c.DeleteHostData(dataKey(host))
}

func dataKey(apiEndpoint string) string {
	return iloAPIKeyPrefix + apiEndpoint
}
