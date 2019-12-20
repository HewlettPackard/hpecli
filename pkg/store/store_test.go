// This is a modified copy of the skv test file.
// Original file is copyright (c) 2016 RapidLoop and
// released under the  MIT licensed at available at:
// https://github.com/rapidloop/skv/blob/master/skv_test.go
// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

const aValue = "this.is.a.value"

func TestKeystoreLocation(t *testing.T) {
	if keystore != "" {
		t.Fatal("keystore shouldn't have a value until Open as been called")
	}
	db, err := NewStore(SKV)
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	if keystore == "" {
		t.Fatal("keystore should have a value after Open as been called")
	}
}

func TestBasic(t *testing.T) {
	db := openForTest(t)
	defer db.Close()
	// put a key
	if err := db.Put("key1", "value1"); err != nil {
		t.Fatal(err)
	}
	// get it back
	var val string
	if err := db.Get("key1", &val); err != nil {
		t.Fatal(err)
	} else if val != "value1" {
		t.Fatalf("got \"%s\", expected \"value1\"", val)
	}
	// put it again with same value
	if err := db.Put("key1", "value1"); err != nil {
		t.Fatal(err)
	}
	// get it back again
	if err := db.Get("key1", &val); err != nil {
		t.Fatal(err)
	} else if val != "value1" {
		t.Fatalf("got \"%s\", expected \"value1\"", val)
	}
	// get something we know is not there
	if err := db.Get("no.such.key", &val); err != ErrNotFound {
		t.Fatalf("got \"%s\", expected absence", val)
	}
	// delete our key
	if err := db.Delete("key1"); err != nil {
		t.Fatal(err)
	}
	// delete it again
	if err := db.Delete("key1"); err != ErrNotFound {
		t.Fatalf("delete returned %v, expected ErrNotFound", err)
	}
	// done
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestMoreNotFoundCases(t *testing.T) {
	db := openForTest(t)
	defer db.Close()
	var val string
	if err := db.Get("key1", &val); err != ErrNotFound {
		t.Fatal(err)
	}
	if err := db.Put("key1", "value1"); err != nil {
		t.Fatal(err)
	}
	if err := db.Delete("key1"); err != nil {
		t.Fatal(err)
	}
	if err := db.Get("key1", &val); err != ErrNotFound {
		t.Fatal(err)
	}
	if err := db.Get("", &val); err != ErrNotFound {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}

type aStruct struct {
	Numbers *[]int
}

func TestRichTypes(t *testing.T) {
	var inval1 = map[string]string{
		"100 meters": "Florence GRIFFITH-JOYNER",
		"200 meters": "Florence GRIFFITH-JOYNER",
		"400 meters": "Marie-José PÉREC",
		"800 meters": "Nadezhda OLIZARENKO",
	}
	var outval1 = make(map[string]string)
	testGetPut(t, inval1, &outval1)
	var inval2 = aStruct{
		Numbers: &[]int{100, 200, 400, 800},
	}
	var outval2 aStruct
	testGetPut(t, inval2, &outval2)
}

func testGetPut(t *testing.T, inval interface{}, outval interface{}) {
	db := openForTest(t)
	defer db.Close()
	input, err := json.Marshal(inval)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Put("test.key", inval); err != nil {
		t.Fatal(err)
	}
	if err := db.Get("test.key", outval); err != nil {
		t.Fatal(err)
	}
	output, err := json.Marshal(outval)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(input, output) {
		t.Fatal("differences encountered")
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestNil(t *testing.T) {
	db := openForTest(t)
	defer db.Close()
	if err := db.Put("key1", nil); err != ErrBadValue {
		t.Fatalf("got %v, expected ErrBadValue", err)
	}
	if err := db.Put("key1", "value1"); err != nil {
		t.Fatal(err)
	}
	// can Get() into a nil value
	if err := db.Get("key1", nil); err != nil {
		t.Fatal(err)
	}
}

func TestGoroutines(t *testing.T) {
	db := openForTest(t)
	defer db.Close()
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			switch rand.Intn(3) {
			case 0:
				if err := db.Put("key1", "value1"); err != nil {
					t.Fatal(err)
				}
			case 1:
				var val string
				if err := db.Get("key1", &val); err != nil && err != ErrNotFound {
					t.Fatal(err)
				}
			case 2:
				if err := db.Delete("key1"); err != nil && err != ErrNotFound {
					t.Fatal(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkPut(b *testing.B) {
	db := openForBenchmark(b)
	defer db.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := db.Put(fmt.Sprintf("key%d", i), aValue); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}

func BenchmarkPutGet(b *testing.B) {
	db := openForBenchmark(b)
	defer db.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := db.Put(fmt.Sprintf("key%d", i), aValue); err != nil {
			b.Fatal(err)
		}
	}
	for i := 0; i < b.N; i++ {
		var val string
		if err := db.Get(fmt.Sprintf("key%d", i), &val); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}

func BenchmarkPutDelete(b *testing.B) {
	db := openForBenchmark(b)
	defer db.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := db.Put(fmt.Sprintf("key%d", i), aValue); err != nil {
			b.Fatal(err)
		}
	}
	for i := 0; i < b.N; i++ {
		if err := db.Delete(fmt.Sprintf("key%d", i)); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}

func openForTest(t *testing.T) Store {
	os.Remove("skv-bench.db")
	db, err := NewStore(SKV)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func openForBenchmark(b *testing.B) Store {
	os.Remove("skv-bench.db")
	db, err := NewStore(SKV)
	if err != nil {
		b.Fatal(err)
	}
	return db
}
