// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package password

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestTwoPasswordOptions(t *testing.T) {
	var p string = "have one"
	err := Read(&p, true, "ilo")

	if err == nil {
		t.Error("expected error")
	}

	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Error("unexpected error text")
	}
}

func TestPasswordAlreadySpecified(t *testing.T) {
	var p string = "have one"
	got := p
	err := Read(&p, false, "ilo")

	if err != nil {
		t.Error("expected error")
	}

	if got != p {
		t.Error("error: password unexpectedly modified")
	}
}

func TestReadFromConsoleWithError(t *testing.T) {
	out := &bytes.Buffer{}
	stdout = out
	reader = &fakePasswordReader{ReturnError: true}

	want := "myPrompt:"
	_, err := ReadFromConsole(want)

	if err == nil {
		t.Fatal("should have seen an error")
	}

	got := out.String()
	if got != want {
		t.Fatalf("didn't get expected prompt written on out.  got=%s, want=%s", got, want)
	}
}

func TestReadFromConsole(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	stdin = in
	stdout = out
	reader = &fakePasswordReader{Password: "mypassword"}

	got, err := ReadFromConsole("abc")
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	if out.String() != "abc\n" {
		t.Error("output should have ended with a \\n")
	}

	if got != "mypassword" {
		t.Fatal("didn't read passsword correctly")
	}
}

func TestReadFromStdin(t *testing.T) {
	in := bytes.NewBufferString(" mypassword \n ")
	stdin = in

	got, err := ReadFromStdIn()
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	// ensure password is read.. and trimmed
	if got != "mypassword" {
		t.Fatal("didn't read passsword correctly")
	}
}

func TestReadFromStdinOnError(t *testing.T) {
	// not putting a \n will cause EOF on read
	in := bytes.NewBufferString("mypassword")
	stdin = in

	_, err := ReadFromStdIn()
	if err == nil {
		t.Fatal("expected EOF error on read without \\n")
	}
}

func TestReadPromptFor(t *testing.T) {
	in := bytes.NewBufferString(" otherPassword \n ")
	stdin = in

	var got string
	err := Read(&got, true, "")

	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	// ensure password is read.. and trimmed
	if got != "otherPassword" {
		t.Fatal("didn't read passsword correctly")
	}
}

func TestReadPrompt(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	stdin = in
	stdout = out
	reader = &fakePasswordReader{Password: "diff-password"}

	wantPrompt := "enter password: "

	var got string
	err := Read(&got, false, wantPrompt)

	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	if out.String() != wantPrompt+"\n" {
		t.Error("output should have ended with a \\n")
	}

	if got != "diff-password" {
		t.Fatal("didn't read passsword correctly")
	}
}

type fakePasswordReader struct {
	Password    string
	ReturnError bool
}

func (pr fakePasswordReader) readPassword() (string, error) {
	if pr.ReturnError {
		return "", errors.New("expected stubbed error")
	}

	return pr.Password, nil
}
