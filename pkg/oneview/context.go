package oneview

import (
	"github.com/HewlettPackard/hpecli/pkg/context"
	"github.com/HewlettPackard/hpecli/pkg/db"
)

func ovContext() (context.Context, error) {
	c, err := context.New(oneViewContextKey, oneViewAPIKeyPrefix, db.Open)
	if err != nil {
		return nil, err
	}

	return c, nil
}
