// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"github.com/HewlettPackard/hpecli/pkg/context"
	"github.com/HewlettPackard/hpecli/pkg/db"
)

const iloAPIKeyPrefix = "hpecli_ilo_token_"
const iloContextKey = "hpecli_ilo_context"

func iloContext() context.Context {
	return context.New(iloContextKey, iloAPIKeyPrefix, db.Open)
}
