// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package store

import "github.com/rapidloop/skv"

type skvstore struct {
	db *skv.KVStore
}

func openSKV(filename string) (Store, error) {
	db, err := skv.Open(filename)
	if err != nil {
		return nil, err
	}

	return &skvstore{db: db}, nil
}

func (s skvstore) Get(key string, value interface{}) error {
	if err := s.db.Get(key, value); err == skv.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	return nil
}

func (s skvstore) Put(key string, value interface{}) error {
	if value == nil {
		return ErrBadValue
	}

	return s.db.Put(key, value)
}

func (s skvstore) Delete(key string) error {
	err := s.db.Delete(key)
	if err == skv.ErrNotFound {
		return ErrNotFound
	}

	return err
}

func (s skvstore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}

	return nil
}
