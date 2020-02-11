package context

import (
	"errors"
	"fmt"

	"github.com/HewlettPackard/hpecli/pkg/db"
)

type Context interface {
	APIKey() (host, apiKey string, err error)
	SetAPIKey(host, sessionKey string) error
	SetContext(host string) error
}

type DBOpen func() (db.Store, error)

type APIContext struct {
	ContextKey   string
	APIKeyPrefix string
	DBOpen       DBOpen
}

var (
	ErrorContextNotFound = errors.New("unable to find specified context")
	ErrorKeyNotFound     = errors.New("unable to find key for specified host")
)

func New(contextKey, apiKeyPrefix string, dbOpen DBOpen) (Context, error) {
	return &APIContext{
		ContextKey:   contextKey,
		APIKeyPrefix: apiKeyPrefix,
		DBOpen:       dbOpen,
	}, nil
}

func (c APIContext) APIKey() (host, sessionKey string, err error) {
	d, err := c.DBOpen()
	if err != nil {
		return "", "", err
	}
	defer d.Close()

	if err := d.Get(c.ContextKey, &host); err != nil {
		return "", "", ErrorContextNotFound
	}

	apiKey := makeAPIKey(c.APIKeyPrefix, host)

	if err := d.Get(apiKey, &sessionKey); err != nil {
		return host, "", ErrorKeyNotFound
	}

	return host, sessionKey, nil
}

func (c APIContext) SetAPIKey(host, sessionKey string) error {
	d, err := c.DBOpen()
	if err != nil {
		return err
	}
	defer d.Close()

	// Save context key
	if e := d.Put(c.ContextKey, host); e != nil {
		return fmt.Errorf("unable to save current context because of %#v", e)
	}

	// Save API Key
	apiKey := makeAPIKey(c.APIKeyPrefix, host)
	if err := d.Put(apiKey, sessionKey); err != nil {
		return fmt.Errorf("unable to save apiKey for %s because of %#v", host, err)
	}

	return nil
}

func (c APIContext) SetContext(host string) error {
	d, err := c.DBOpen()
	if err != nil {
		return err
	}
	defer d.Close()

	// Save context key
	if err := d.Put(c.ContextKey, host); err != nil {
		return fmt.Errorf("unable to save current context because of %#v", err)
	}

	return nil
}

func makeAPIKey(apiKeyPrefix, host string) string {
	return fmt.Sprintf("%s-%s", apiKeyPrefix, host)
}
