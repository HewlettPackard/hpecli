package password

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// These are defined here for testing and should be
// needed outside of that
var (
	stdout io.Writer
	stdin  io.Reader
	reader passwordReader
)

// interface to allow us to test reading from the console
type passwordReader interface {
	readPassword() (string, error)
}

type stdInPasswordReader struct{}

func init() {
	stdout = os.Stdout
	stdin = os.Stdin
	reader = &stdInPasswordReader{}
}

// readPassword reads password from stdin
func (pr stdInPasswordReader) readPassword() (string, error) {
	//nolint:unconvert  // false positive lint err
	pwd, err := terminal.ReadPassword(int(syscall.Stdin))

	return string(pwd), err
}

// ReadFromConsole write the prompt to the console and then
// will read the password from the user console without
// echoing the password to the console
func ReadFromConsole(prompt string) (string, error) {
	fmt.Fprint(stdout, prompt)

	p, err := reader.readPassword()
	if err != nil {
		return "", err
	}

	// rather than force everyone to do a linefeed after reading
	// we'll just do it here
	fmt.Fprint(stdout, "\n")

	return p, nil
}

// ReadFromStdIn reads from os.Stdin and trims any CR & LF characters
// as well as any spaces around the input.  This method allows passwords
// to be passed to the CLI with commands like:
//   <Windows>
//     echo mypassword | hpecli ilo login --host 1.1.1.1 -u admin --password-stdin
func ReadFromStdIn() (string, error) {
	reader := bufio.NewReader(stdin)

	t, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading from stdin: %w", err)
	}

	t = strings.TrimSpace(t)

	return t, nil
}


// Read will use one of several methods to get a password from a user
// password is an existing password if set.  This is return unmodified if
// passed in with a value.
// readStdIn is true if the password should be read from standard in
// prompt is what the console prompt should be if it should be read
// from the console from the user
func Read(password *string, readStdIn bool, prompt string) (err error) {
	if *password != "" && readStdIn {
		return errors.New("--password and --password-stdin are mutually exclusive")
	}

	if *password != "" {
		// password specified.. leave it and don't get another one
		return
	}

	// asked to read from stdin
	if readStdIn {
		*password, err = ReadFromStdIn()
		return
	}

	// prompt the user for the password on the console
	*password, err = ReadFromConsole(prompt)

	return
}
