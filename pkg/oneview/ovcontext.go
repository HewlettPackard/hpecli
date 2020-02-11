package oneview

import (
	"github.com/HewlettPackard/hpecli/pkg/context"
	"github.com/HewlettPackard/hpecli/pkg/db"
)

const oneViewAPIKeyPrefix = "hpecli_oneview_token_"
const oneViewContextKey = "hpecli_oneview_context"

func ovContext() (context.Context, error) {
	c, err := context.New(oneViewContextKey, oneViewAPIKeyPrefix, db.Open)
	if err != nil {
		return nil, err
	}

	return c, nil
}
