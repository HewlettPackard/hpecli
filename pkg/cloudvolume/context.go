// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package cloudvolume

import (
	"github.com/HewlettPackard/hpecli/pkg/context"
	"github.com/HewlettPackard/hpecli/pkg/db"
)

const cvAPIKeyPrefix = "hpecli_cloudvolume_token_"
const cvContextKey = "hpecli_cloudvolume_context"

func cvContext() context.Context {
	return context.NewContext(cvContextKey, cvAPIKeyPrefix, db.Open)
}
