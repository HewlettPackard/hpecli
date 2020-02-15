package context

import (
	"fmt"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/db"
)

const (
	apiKeyPrefix  = "somePrefix"
	contextKey    = "someContext"
	host          = "someHost"
	key           = "someKey"
	errTempl      = "got: %s, wanted: %s"
	fail          = "fail"
	errExpected   = "error was expected"
	errUnexpected = "unexpected error: %v"
)

func TestNewContext(t *testing.T) {
	c := New(contextKey, apiKeyPrefix, MockOpen)

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
		t.Fatalf(errUnexpected, err)
	}
}

func TestGetAPIKey(t *testing.T) {
	c := withMockStore()

	if err := c.SetAPIKey(host, key); err != nil {
		t.Fatalf(errUnexpected, err)
	}

	var got string
	err := c.APIKey(&got)
	if err != nil {
		t.Fatal(err)
	}

	if got != key {
		t.Fatal("wrong host returned")
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

	var v string
	err := c.APIKey(&v)
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

	var key string
	err = c.APIKey(&key)
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

	err := c.APIKey(nil)
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

func TestChangeContextPutFails(t *testing.T) {
	c := withMockStore()

	if err := c.ChangeContext("fail"); err == nil {
		t.Fatalf(errExpected)
	}
}

func TestChangeContextWritesValue(t *testing.T) {
	c := withMockStore()

	if err := c.ChangeContext(host); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got string

	err := ms.Get(contextKey, &got)
	if err != nil {
		t.Fatal(err)
	}

	if got != host {
		t.Fatal("didn't retrieve expected value after ChangeContext")
	}
}

func FailOpen() (db.Store, error) {
	return nil, fmt.Errorf(errExpected)
}

func withMockStore() Context {
	var c Context = &APIContext{
		APIKeyPrefix: apiKeyPrefix,
		ContextKey:   contextKey,
		DBOpen:       MockOpen,
	}

	return c
}
