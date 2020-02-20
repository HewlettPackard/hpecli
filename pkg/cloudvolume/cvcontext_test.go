package cloudvolume

import (
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestStoreContext(t *testing.T) {
	if err := storeContext("host1", "blahKey"); err != nil {
		t.Fatal(err)
	}
}

func TestStoreContextFailsEmptyKey(t *testing.T) {
	if err := storeContext("", "blahKey"); err == nil {
		t.Fatal("expected failure no empty key")
	}
}

func TestStoreGetWorks(t *testing.T) {
	const h1 = "host1"

	const v1 = "value1"

	if err := storeContext(h1, v1); err != nil {
		t.Fatal(err)
	}

	d, err := getContext()
	if err != nil {
		t.Fatal(err)
	}

	if d.Host != h1 {
		t.Fatal("didn't retrieve matching host value")
	}

	if d.APIKey != v1 {
		t.Fatal("didn't retrieve matching apikey value")
	}
}
