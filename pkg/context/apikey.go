package context

import (
	"errors"
	"fmt"

	"github.com/HewlettPackard/hpecli/pkg/db"
)

type Context interface {
	APIKey(value interface{}) error
	SetAPIKey(key string, value interface{}) error
	ChangeContext(key string) error
}

type DBOpenFunc func() (db.Store, error)

var DefaultDBOpenFunc = db.Open

type APIContext struct {
	ContextKey   string
	APIKeyPrefix string
	DBOpen       DBOpenFunc
}

var (
	ErrorContextNotFound = errors.New("unable to find specified context")
	ErrorKeyNotFound     = errors.New("unable to find key for specified host")
	ErrorInvalidKey      = errors.New("invalid key specified.  Key can not be empty")
)

func NewWithDB(contextKey, apiKeyPrefix string, dbOpen DBOpenFunc) Context {
	return &APIContext{
		ContextKey:   contextKey,
		APIKeyPrefix: apiKeyPrefix,
		DBOpen:       dbOpen,
	}
}

func New(contextKey, apiKeyPrefix string) Context {
	return NewWithDB(contextKey, apiKeyPrefix, DefaultDBOpenFunc)
}

func (c APIContext) APIKey(value interface{}) error {
	d, err := c.DBOpen()
	if err != nil {
		return err
	}
	defer d.Close()

	var host string
	if err := d.Get(c.ContextKey, &host); err != nil {
		return ErrorContextNotFound
	}

	apiKey := makeAPIKey(c.APIKeyPrefix, host)

	if err := d.Get(apiKey, value); err != nil {
		return ErrorKeyNotFound
	}

	return nil
}

func (c APIContext) SetAPIKey(key string, value interface{}) error {
	if key == "" {
		return ErrorInvalidKey
	}

	d, err := c.DBOpen()
	if err != nil {
		return err
	}
	defer d.Close()

	// Save context key
	if e := d.Put(c.ContextKey, key); e != nil {
		return fmt.Errorf("unable to save current context because of %#v", e)
	}

	// Save API Key
	apiKey := makeAPIKey(c.APIKeyPrefix, key)
	if err := d.Put(apiKey, value); err != nil {
		return fmt.Errorf("unable to save apiKey for %s because of %#v", key, err)
	}

	return nil
}

func (c APIContext) ChangeContext(key string) error {
	if key == "" {
		return ErrorInvalidKey
	}

	d, err := c.DBOpen()
	if err != nil {
		return err
	}
	defer d.Close()

	// Save context key
	if err := d.Put(c.ContextKey, key); err != nil {
		return fmt.Errorf("unable to save current context because of %#v", err)
	}

	return nil
}

func makeAPIKey(apiKeyPrefix, host string) string {
	return apiKeyPrefix + host
}
