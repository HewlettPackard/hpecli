package context

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/db"
)

type MockStore struct {
	m map[string][]byte
}

var ms = &MockStore{m: make(map[string][]byte)}

func MockOpen() (db.Store, error) {
	return ms, nil
}

func (ms *MockStore) Get(key string, value interface{}) error {
	if strings.Contains(key, "fail") {
		return errors.New("expected error")
	}

	d := gob.NewDecoder(bytes.NewReader(ms.m[key]))

	return d.Decode(value)
}

func (ms *MockStore) Put(key string, value interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}

	ms.m[key] = buf.Bytes()
	v := fmt.Sprintf("%s", value)

	if strings.Contains(key, "fail") || strings.Contains(v, "fail") {
		return errors.New("expected error")
	}

	return nil
}

func (ms *MockStore) Delete(key string) error {
	if strings.Contains(key, "fail") {
		return errors.New("expected error")
	}
	delete(ms.m, key)

	return nil
}

func (ms *MockStore) Close() error {
	return nil
}

func (ms *MockStore) Clear() error {
	ms.m = make(map[string][]byte)
	return nil
}
