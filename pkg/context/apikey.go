package context

import (
	"fmt"

	"github.com/HewlettPackard/hpecli/pkg/db"
)

type Context interface {
	APIKey() (host, apiKey string, err error)
	SetAPIKey(host, sessionKey string) error
}

type DBOpen func() (db.Store, error)

type APIContext struct {
	ContextKey   string
	APIKeyPrefix string
	DBOpen       DBOpen
}

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
		err = fmt.Errorf("unable to retrieve current context: %#v", err)
		return "", "", err
	}

	apiKey := makeAPIKey(c.APIKeyPrefix, host)

	if err := d.Get(apiKey, &sessionKey); err != nil {
		err = fmt.Errorf("unable to retrieve current context: %#v", err)
		return "", "", err
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
		return fmt.Errorf("unable to retrieve current context because of %#v", e)
	}

	// Save API Key
	apiKey := makeAPIKey(c.APIKeyPrefix, host)
	if e := d.Put(apiKey, sessionKey); e != nil {
		return fmt.Errorf("unable to save apiKey for %s because of %#v", host, e)
	}

	return nil
}

func makeAPIKey(apiKeyPrefix, host string) string {
	return fmt.Sprintf("%s-%s", apiKeyPrefix, host)
}
