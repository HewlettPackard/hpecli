package cloudvolume

import (
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/context"
)

func TestContextNotNil(t *testing.T) {
	c := cvContext()
	if c == nil {
		t.Fatal("cvContext() returned nil")
	}

	if cvAPIKeyPrefix != c.(*context.APIContext).APIKeyPrefix {
		t.Fatalf("APIKeyPrefix wasn't stored correctly")
	}

	if cvContextKey != c.(*context.APIContext).ContextKey {
		t.Fatalf("ContextKey wasn't stored correctly")
	}
}
