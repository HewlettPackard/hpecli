package db

import (
	"io"
	"os"
	"testing"
)

const fname = "TEST_FILE"

func TestOpenMissingFile(t *testing.T) {
	_, err := openSKV("")
	if err == nil {
		t.Fatal("open should fail with empty filename")
	}
}

func TestOpenSKV(t *testing.T) {
	db, err := openSKV(fname)
	defer cleanupSKV(db, fname)

	if err != nil {
		t.Fatal("expected to be able to open, but failed")
	}

	if db == nil {
		t.Fatal("db was opened but is nil")
	}
}

func TestBadGet(t *testing.T) {
	db, _ := openSKV(fname)
	defer cleanupSKV(db, fname)

	if err := db.Put("key1", "value1"); err != nil {
		t.Fatal(err)
	}

	var val int
	// put a string, now try and retrieve as int
	if err := db.Get("key1", &val); err == nil {
		t.Fatal("didn't report error on bad get")
	}
}

func cleanupSKV(db io.Closer, f string) {
	if db != nil {
		db.Close()
	}

	os.Remove(f)
}
