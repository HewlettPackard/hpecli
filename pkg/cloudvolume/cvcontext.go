// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"github.com/HewlettPackard/hpecli/pkg/context"
	"github.com/HewlettPackard/hpecli/pkg/db"
)

const cvAPIKeyPrefix = "hpecli_cloudvolume_token_"
const cvContextKey = "hpecli_cloudvolume_context"

type cvContextData struct {
	Host   string
	APIKey string
}

func storeContext(key, token string) error {
	c := context.New(cvContextKey, cvAPIKeyPrefix, db.Open)
	return c.SetAPIKey(key, &cvContextData{key, token})
}

func getContext() (*cvContextData, error) {
	c := context.New(cvContextKey, cvAPIKeyPrefix, db.Open)
	var d cvContextData
	err := c.APIKey(&d)
	return &d, err
}
