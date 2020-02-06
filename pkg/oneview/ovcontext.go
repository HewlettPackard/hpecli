package oneview

import (
	"github.com/HewlettPackard/hpecli/pkg/context"
	"github.com/HewlettPackard/hpecli/pkg/store"
)

func ovContext() (context.Context, error) {
	c, err := context.NewContext(oneViewContextKey, oneViewAPIKeyPrefix, store.Open)
	if err != nil {
		return nil, err
	}

	return c, nil
}
