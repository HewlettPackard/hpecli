package context

import (
	"fmt"

	"github.com/HewlettPackard/hpecli/pkg/store"
)

type Context interface {
	APIKey() (host, apiKey string, err error)
	SetAPIKey(host, sessionKey string) error
}

type DBOpen func() (store.Store, error)

type APIContext struct {
	ContextKey   string
	APIKeyPrefix string
	DBOpen       DBOpen
}

func NewContext(contextKey, apiKeyPrefix string, dbOpen DBOpen) (Context, error) {
	return &APIContext{
		ContextKey:   contextKey,
		APIKeyPrefix: apiKeyPrefix,
		DBOpen:       dbOpen,
	}, nil
}

func (c APIContext) APIKey() (host, sessionKey string, err error) {
	db, err := c.DBOpen()
	if err != nil {
		return "", "", err
	}
	defer db.Close()

	if err := db.Get(c.ContextKey, &host); err != nil {
		err := fmt.Errorf("Unable to retrieve current context: %#v", err)
		return "", "", err
	}

	apiKey := makeAPIKey(c.APIKeyPrefix, host)

	if err := db.Get(apiKey, &sessionKey); err != nil {
		err := fmt.Errorf("Unable to retrieve current context: %#v", err)
		return "", "", err
	}

	return host, sessionKey, nil
}

func (c APIContext) SetAPIKey(host, sessionKey string) error {
	db, err := c.DBOpen()
	if err != nil {
		return err
	}
	defer db.Close()

	// Save context key
	if e := db.Put(c.ContextKey, host); e != nil {
		return fmt.Errorf("unable to retrieve current context because of %#v", e)
	}

	// Save API Key
	apiKey := makeAPIKey(c.APIKeyPrefix, host)
	if e := db.Put(apiKey, sessionKey); e != nil {
		return fmt.Errorf("unable to save apiKey for %s because of %#v", host, e)
	}

	return nil
}

func makeAPIKey(apiKeyPrefix, host string) string {
	return fmt.Sprintf("%s-%s", apiKeyPrefix, host)
}
