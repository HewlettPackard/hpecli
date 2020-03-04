package autocomplete

import (
	"testing"
)

func TestCompletion_bash(t *testing.T) {
	err := run(Cmd, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompletion_not_bash(t *testing.T) {
	genautocompleteCmd.autocompleteType = "zsh"
	err := run(Cmd, nil)

	if err != nil {
		t.Fatal(err)
	}
}
