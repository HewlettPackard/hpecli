package cloudvolume

import (
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/context"
)

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestSaveData(t *testing.T) {
	if err := saveData("host1.cloud", "blah-token-value"); err != nil {
		t.Fatal(err)
	}
}

func TestGetWorks(t *testing.T) {
	const h1 = "host1"

	const v1 = "value1"

	if err := saveData(h1, v1); err != nil {
		t.Fatal(err)
	}

	gotHost, gotToken, err := hostAndToken()
	if err != nil {
		t.Fatal(err)
	}

	if gotHost != h1 {
		t.Fatal("didn't retrieve matching host")
	}

	if gotToken != v1 {
		t.Fatal("didn't retrieve matching value")
	}
}
