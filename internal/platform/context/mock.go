package context

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/db"
)

var ErrorExpected = errors.New("expected error")

type MockStore struct {
	m map[string][]byte
}

var ms = &MockStore{m: make(map[string][]byte)}

func MockOpen() (db.Store, error) {
	return ms, nil
}

func (ms *MockStore) Get(key string, value interface{}) error {
	if strings.Contains(key, "fail") {
		return ErrorExpected
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
		return ErrorExpected
	}

	return nil
}

func (ms *MockStore) Delete(key string) error {
	if strings.Contains(key, "fail") {
		return ErrorExpected
	}

	if _, ok := ms.m[key]; !ok {
		return db.ErrNotFound
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
