// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package db

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// StorageEngine defines the backend storage engine to be used
// This ostensibly allows different backends to be used.
// e.g. if we want to switch to a json based file
type StorageEngine int

const (
	// store.db will get stored in a directory called .hpe in the
	// users home directory
	filename = "store.db"
	hpeDir   = ".hpe"
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

	keystore = KeystoreLocation()
)

// Open is what should generally be used to get access to the
// store.  There is currently a single backend - so no need to specify
func Open() (Store, error) {
	return NewStore(SKV)
}

// NewStore returns a handle to a store.  Currently there is a single
// backend - so you can use DefautlStore instead of this method
func NewStore(se StorageEngine) (Store, error) {
	ensureDirExists(keystore)

	if se == SKV {
		return openSKV(keystore)
	}

	return nil, errors.New("unknown StorageEngine type specified")
}

// KeystoreLocation returns the file path to store the DB for individual
// users.  This will be ~/.hpe/store.db
func KeystoreLocation() string {
	// in case we can't get userhomedir - we will use just the filename
	// which will then use the system default path
	result := filename

	if userDir, err := os.UserHomeDir(); err != nil {
		logrus.Warningf("Unable to retrieve users home directory: %v", err)
	} else {
		result = filepath.ToSlash(path.Join(userDir, hpeDir, filename))
	}

	return result
}

func ensureDirExists(keystore string) {
	dir := path.Dir(keystore)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		logrus.Warningf("Unable to persist data in \"%s\" because of error: %v", dir, err)
	}
}
