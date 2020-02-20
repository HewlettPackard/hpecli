// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package oneview

import (
	"github.com/HewlettPackard/hpecli/pkg/context"
)

const oneViewAPIKeyPrefix = "hpecli_oneview_token_"
const oneViewContextKey = "hpecli_oneview_context"

type ovContextData struct {
	Host   string
	APIKey string
}

func storeContext(key, token string) error {
	c := context.New(oneViewContextKey, oneViewAPIKeyPrefix)
	return c.SetAPIKey(key, &ovContextData{key, token})
}

func getContext() (*ovContextData, error) {
	c := context.New(oneViewContextKey, oneViewAPIKeyPrefix)

	var d ovContextData
	err := c.APIKey(&d)

	return &d, err
}

func changeContext(key string) error {
	c := context.New(oneViewContextKey, oneViewAPIKeyPrefix)
	return c.ChangeContext(key)
}
