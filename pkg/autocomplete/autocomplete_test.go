package autocomplete

import (
	"testing"
)

func TestCompletionForBash(t *testing.T) {
	err := run(Cmd, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompletionForNotBash(t *testing.T) {
	genautocompleteCmd.autocompleteType = "zsh"
	err := run(Cmd, nil)

	if err != nil {
		t.Fatal(err)
	}
}
