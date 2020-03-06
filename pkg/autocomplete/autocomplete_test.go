package autocomplete

import (
	"testing"
)

func TestCompletionForBash(t *testing.T) {
	cmd := NewAutoCompleteCommand()

	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompletionForNotBash(t *testing.T) {
	opts := &autocompleteOptions{
		acType: "zsh",
	}

	err := runAutoComplete(nil, opts)
	if err == nil {
		t.Fatal("expected to find error on non bash shell type")
	}
}
