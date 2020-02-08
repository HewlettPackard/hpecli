package context

import (
	"fmt"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/db"
)

const (
	apiKeyPrefix = "somePrefix"
	contextKey   = "someContext"
	host         = "someHost"
	key          = "someKey"
	errTempl     = "got: %s, wanted: %s"
	fail         = "fail"
	errExpected  = "error was expected"
)

func TestNewContext(t *testing.T) {
	c, err := New(contextKey, apiKeyPrefix, MockOpen)
	if err != nil {
		t.Fatal(err)
	}

	if c.(*APIContext).ContextKey != contextKey {
		t.Fatal("ContextKey value not set as expected")
	}

	if c.(*APIContext).APIKeyPrefix != apiKeyPrefix {
		t.Fatal("ContextKey value not set as expected")
	}
}

func TestSetAPIKey(t *testing.T) {
	c := withMockStore()

	if err := c.SetAPIKey(host, key); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetAPIKey(t *testing.T) {
	c := withMockStore()

	if err := c.SetAPIKey(host, key); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got1, got2, err := c.APIKey()
	if err != nil {
		t.Fatal(err)
	}

	if got1 != host {
		t.Fatal("wrong host returned")
	}

	if got2 != key {
		t.Fatal("wrong key returned")
	}
}

func TestSetContextFails(t *testing.T) {
	c := withMockStore()
	c.(*APIContext).ContextKey = fail

	err := c.SetAPIKey(host, key)
	if err == nil {
		t.Fatal(errExpected)
	}
}

func TestGetContextFails(t *testing.T) {
	c := withMockStore()
	c.(*APIContext).ContextKey = fail

	_, _, err := c.APIKey()
	if err == nil {
		t.Fatal(errExpected)
	}
}

func TestSetSessionKeyFails(t *testing.T) {
	c := withMockStore()
	err := c.SetAPIKey(host, fail)

	if err == nil {
		t.Fatal(errExpected)
	}
}

func TestWhenAPIKeyFails(t *testing.T) {
	d, err := MockOpen()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	s := d.(*MockStore)
	s.Put(contextKey, fail)

	c := withMockStore()

	_, _, err = c.APIKey()
	if err == nil {
		t.Fatal(errExpected)
	}

	err = c.SetAPIKey(host, fail)
	if err == nil {
		t.Fatal(errExpected)
	}
}

func TestMakeAPIKey(t *testing.T) {
	got := makeAPIKey(apiKeyPrefix, "")
	if !strings.HasPrefix(got, apiKeyPrefix) {
		t.Fatalf(errTempl, got, apiKeyPrefix)
	}

	got = makeAPIKey("", host)
	if !strings.HasSuffix(got, host) {
		t.Fatalf(errTempl, got, apiKeyPrefix)
	}

	got = makeAPIKey(apiKeyPrefix, host)
	if !strings.HasPrefix(got, apiKeyPrefix) {
		t.Fatalf(errTempl, got, apiKeyPrefix+host)
	}

	if !strings.HasSuffix(got, host) {
		t.Fatalf(errTempl, got, apiKeyPrefix+host)
	}
}

func TestGetAPIFailsOnDBOpen(t *testing.T) {
	var c Context = &APIContext{
		APIKeyPrefix: apiKeyPrefix,
		ContextKey:   contextKey,
		DBOpen:       FailOpen,
	}

	_, _, err := c.APIKey()
	if err == nil {
		t.Fatal(errExpected)
	}
}

func TestSetAPIFailsOnDBOpen(t *testing.T) {
	var c Context = &APIContext{
		APIKeyPrefix: apiKeyPrefix,
		ContextKey:   contextKey,
		DBOpen:       FailOpen,
	}

	err := c.SetAPIKey(host, apiKeyPrefix)
	if err == nil {
		t.Fatal(errExpected)
	}
}

func FailOpen() (db.Store, error) {
	return nil, fmt.Errorf("expected")
}

func withMockStore() Context {
	var c Context = &APIContext{
		APIKeyPrefix: apiKeyPrefix,
		ContextKey:   contextKey,
		DBOpen:       MockOpen,
	}

	return c
}
