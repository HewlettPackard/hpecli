// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package store

import (
	"errors"
	"os"
	"path"

	"github.com/HewlettPackard/hpecli/pkg/logger"
)

// StorageEngine defines the backend storage engine to be used
// This ostensibly allows different backends to be used.
// e.g. if we want to switch to a json based file
type StorageEngine int

const (
	filename = "hpecli_store.db"
	// SKV is Simple Key Value store from: github.com/rapidloop/skv
	SKV StorageEngine = iota
)

// Store is the interface to store persistent data
// This allows for changing backend storage engines
type Store interface {
	// Get an entry from the store. "value" must be a pointer-typed. If the key
	// is not present in the store, Get returns ErrNotFound.
	//
	//	type MyStruct struct {
	//	    Numbers []int
	//	}
	//	var val MyStruct
	//	if err := store.Get("key42", &val); err == store.ErrNotFound {
	//	    // "key42" not found
	//	} else if err != nil {
	//	    // an error occurred
	//	} else {
	//	    // ok
	//	}
	//
	// The value passed to Get() can be nil, in which case any value read from
	// the store is silently discarded.
	//
	//  if err := store.Get("key42", nil); err == nil {
	//      fmt.Println("entry is present")
	//  }
	Get(key string, value interface{}) error
	// Put an entry into the store. The passed value is gob-encoded and stored.
	// The key can be an empty string, but the value cannot be nil - if it is,
	// Put() returns ErrBadValue.
	//
	//	err := store.Put("key42", 156)
	//	err := store.Put("key42", "this is a string")
	//	m := map[string]int{
	//	    "harry": 100,
	//	    "emma":  101,
	//	}
	//	err := store.Put("key43", m)
	Put(key string, value interface{}) error
	// Delete the entry with the given key. If no such key is present in the store,
	// it returns ErrNotFound.
	//
	//	store.Delete("key42")
	Delete(key string) error
	// Close closes the key-value store file.
	Close() error
}

var (
	// ErrNotFound is returned when the key supplied to a Get or Delete
	// method does not exist in the database.
	ErrNotFound = errors.New("db: key not found")

	// ErrBadValue is returned when the value supplied to the Put method is nil.
	ErrBadValue = errors.New("db: bad value")

	keystore = ""
)

// Open is what should generaly be used to get access to the
// store.  There is currently a single backend - so no need to specify
func Open() (Store, error) {
	return NewStore(SKV)
}

// NewStore returns a handle to a store.  Currently there is a single
// backend - so you can use DefautlStore instead of this method
func NewStore(st StorageEngine) (Store, error) {
	if keystore == "" {
		keystore = keystoreLocation()
	}
	switch st {
	case SKV:
		return openSKV(keystore)
	}
	return nil, errors.New("unknown StorageEngine type specified")
}

// return the filename in the users home directory
func keystoreLocation() string {
	if dir, err := os.UserHomeDir(); err != nil {
		logger.Warning("Unable to retrieve users home directory: %v", err)
	} else {
		return path.Join(dir, filename)
	}
	return filename
}
