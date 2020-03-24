package autocomplete

import (
	"testing"
)

func TestCompletionForBash(t *testing.T) {
	cmd := NewAutoCompleteCommand()
	cmd.SetArgs([]string{"--type", "bash"})

	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompletionForZsh(t *testing.T) {
	cmd := NewAutoCompleteCommand()
	cmd.SetArgs([]string{"--type", "zsh"})

	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompletionForPowershell(t *testing.T) {
	cmd := NewAutoCompleteCommand()
	cmd.SetArgs([]string{"--type", "powershell"})

	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
func TestCompletionForNotBashZshPowershell(t *testing.T) {
	opts := &autocompleteOptions{
		acType: "pwr",
	}

	err := runAutoComplete(nil, opts)
	if err == nil {
		t.Fatal("expected to find error on non bash shell type")
	}
}

func TestCompletionForInvalidFIlePathPS(t *testing.T) {
	cmd := NewAutoCompleteCommand()
	cmd.SetArgs([]string{"--type", "powershell", "--completionfile", "/"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected to find error on invalid file path")
	}
}

func TestCompletionForInvalidFIlePathBash(t *testing.T) {
	cmd := NewAutoCompleteCommand()
	cmd.SetArgs([]string{"--type", "bash", "--completionfile", "/"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected to find error on invalid file path")
	}
}

func TestCompletionForInvalidFIlePathZShell(t *testing.T) {
	cmd := NewAutoCompleteCommand()
	cmd.SetArgs([]string{"--type", "zsh", "--completionfile", "/"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected to find error on invalid file path")
	}
}
