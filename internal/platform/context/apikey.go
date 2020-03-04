package context

import (
	"errors"
	"fmt"

	"github.com/HewlettPackard/hpecli/internal/platform/db"
)

type Context interface {
	ModuleContext() (string, error)
	SetModuleContext(value string) error
	HostData(hostKey string, value interface{}) error
	SetHostData(hostKey string, data interface{}) error
	DeleteHostData(hostKey string) error
}

type DBOpenFunc func() (db.Store, error)

var DefaultDBOpenFunc = db.Open

type APIContext struct {
	contextKey string
	dbOpenFunc DBOpenFunc
}

var (
	ErrorContextNotFound = errors.New("unable to find specified context")
	ErrorKeyNotFound     = errors.New("unable to find the specified key")
	ErrorInvalidKey      = errors.New("invalid key specified.  Key can not be empty")
	ErrorInvalidValue    = errors.New("invalid value specified.  Value can not be nil")
)

func NewWithDB(contextKey string, dbOpen DBOpenFunc) Context {
	return &APIContext{
		contextKey: contextKey,
		dbOpenFunc: dbOpen,
	}
}

func New(contextKey string) Context {
	return NewWithDB(contextKey, DefaultDBOpenFunc)
}

func (c APIContext) ModuleContext() (string, error) {
	d, err := c.dbOpenFunc()
	if err != nil {
		return "", err
	}

	defer d.Close()

	var value string
	if err := d.Get(c.contextKey, &value); err != nil {
		return "", ErrorContextNotFound
	}

	return value, nil
}

func (c APIContext) SetModuleContext(value string) error {
	d, err := c.dbOpenFunc()
	if err != nil {
		return err
	}

	defer d.Close()

	// Save context key
	if err := d.Put(c.contextKey, value); err != nil {
		return fmt.Errorf("unable to save the context: %w", err)
	}

	return nil
}

func (c APIContext) HostData(hostKey string, value interface{}) error {
	if hostKey == "" {
		return ErrorInvalidKey
	}

	d, err := c.dbOpenFunc()
	if err != nil {
		return err
	}

	defer d.Close()

	if err := d.Get(hostKey, value); err != nil {
		return ErrorKeyNotFound
	}

	return nil
}

func (c APIContext) SetHostData(hostKey string, value interface{}) error {
	if hostKey == "" {
		return ErrorInvalidKey
	}

	if value == nil {
		return ErrorInvalidValue
	}

	d, err := c.dbOpenFunc()
	if err != nil {
		return err
	}

	defer d.Close()

	// Save context key
	if err := d.Put(hostKey, value); err != nil {
		return fmt.Errorf("unable to save host value: %w", err)
	}

	return nil
}

func (c APIContext) DeleteHostData(hostKey string) error {
	if hostKey == "" {
		return ErrorInvalidKey
	}

	d, err := c.dbOpenFunc()
	if err != nil {
		return err
	}

	defer d.Close()

	// delete the data
	if err := d.Delete(hostKey); err != nil {
		return fmt.Errorf("unable to delete host data: %w", err)
	}

	return nil
}
