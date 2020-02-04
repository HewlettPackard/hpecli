package greenlake

import (
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/HewlettPackard/hpecli/pkg/store"
)

const errTempl = "got: %s, wanted: %s"
const errMsg = "default values should be empty"

func TestMakeAPIKey(t *testing.T) {
	got := makeAPIKey("")
	if got != greenLakeAPIKeyPrefix {
		t.Fatalf(errTempl, got, greenLakeAPIKeyPrefix)
	}

	const testHost1 = "SomeFunkyHost"

	got = makeAPIKey(testHost1)
	if !strings.HasPrefix(got, greenLakeAPIKeyPrefix) {
		t.Fatalf(errTempl, got, greenLakeAPIKeyPrefix+testHost1)
	}

	if !strings.HasSuffix(got, testHost1) {
		t.Fatalf(errTempl, got, greenLakeAPIKeyPrefix+testHost1)
	}
}

func TestMakeTenantID(t *testing.T) {
	got := makeTenantID("")
	if got != greenLakeTenantIDPrefix {
		t.Fatalf(errTempl, got, greenLakeTenantIDPrefix)
	}

	const testHost1 = "SomeFunkyHost"

	got = makeTenantID(testHost1)
	if !strings.HasPrefix(got, greenLakeTenantIDPrefix) {
		t.Fatalf(errTempl, got, greenLakeTenantIDPrefix+testHost1)
	}

	if !strings.HasSuffix(got, testHost1) {
		t.Fatalf(errTempl, got, greenLakeTenantIDPrefix+testHost1)
	}
}

func TestSetAPIKey(t *testing.T) {
	const h1 = "someHost1"

	const t1 = "someTenantID"

	const v1 = "valueToStore"

	_ = setTokenTenantID(h1, t1, v1)

	db, err := store.Open()
	if err != nil {
		logger.Debug("unable to open keystore: %v", err)
		return
	}

	defer db.Close()

	apiKey := makeAPIKey(h1)

	var got string
	if err := db.Get(apiKey, &got); err != nil {
		t.Fatal(err)
	}

	if got != v1 {
		t.Fatalf(errTempl, got, v1)
	}
}

func TestGetTokenTenantID(t *testing.T) {
	const h1 = "someHost1"

	const t1 = "someTenantID"

	const v1 = "valueToStore"

	_ = setTokenTenantID(h1, t1, v1)

	gotHost, gotTenantID, gotAPIKey := getTokenTenantID()

	if gotHost != h1 {
		t.Fatalf(errTempl, gotHost, h1)
	}

	if gotTenantID != t1 {
		t.Fatalf(errTempl, gotTenantID, h1)
	}

	if gotAPIKey != v1 {
		t.Fatalf(errTempl, gotAPIKey, v1)
	}
}

func TestGetAPIKeyFailDBOpenReturnsEmptyDefaults(t *testing.T) {
	db, _ := store.Open()
	defer db.Close()

	// fails because db is already open
	host, tenant, key := getTokenTenantID()
	if host != "" && tenant != "" && key != "" {
		t.Fatal(errMsg)
	}
}

func TestDBDoesntHaveContextReturnsEmptyDefaults(t *testing.T) {
	db, err := store.Open()
	if err != nil {
		t.Fatal(err)
	}

	_ = db.Delete(greenLakeContextKey)
	db.Close()

	host, tenant, key := getTokenTenantID()
	if host != "" && tenant != "" && key != "" {
		t.Fatal(errMsg)
	}
}

func TestGetAPIKeyDBDoesntHaveHostReturnsEmptyDefaults(t *testing.T) {
	const h1 = "host1"

	const t1 = "tenant1"

	const v1 = "value1"

	_ = setTokenTenantID(h1, t1, v1)

	db, err := store.Open()
	if err != nil {
		t.Fatal(err)
	}

	_ = db.Delete(makeAPIKey(h1))
	db.Close()

	host, tenant, key := getTokenTenantID()
	if host != "" && tenant != "" && key != "" {
		t.Fatal(errMsg)
	}
}

func TestSetAPIKeyFailWithDBOpen(t *testing.T) {
	db, err := store.Open()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	err = setTokenTenantID("", "", "")
	if err == nil {
		t.Fatal("expected error that db was already open")
	}
}
