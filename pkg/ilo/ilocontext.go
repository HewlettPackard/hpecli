// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"github.com/HewlettPackard/hpecli/pkg/context"
)

const iloAPIKeyPrefix = "hpecli_ilo_token_"
const iloContextKey = "hpecli_ilo_context"

type iloContextData struct {
	Host   string
	APIKey string
}

func storeContext(key, token string) error {
	c := context.New(iloContextKey, iloAPIKeyPrefix)
	return c.SetAPIKey(key, &iloContextData{key, token})
}

func getContext() (*iloContextData, error) {
	c := context.New(iloContextKey, iloAPIKeyPrefix)

	var d iloContextData
	err := c.APIKey(&d)

	return &d, err
}

func changeContext(key string) error {
	c := context.New(iloContextKey, iloAPIKeyPrefix)
	return c.ChangeContext(key)
}
