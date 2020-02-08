package oneview

import (
	"testing"
)

func TestInitContextMethod(t *testing.T) {
	c := initContext(t)
	if c == nil {
		t.Fatal("expected ovContext variable to be initialized")
	}
}
