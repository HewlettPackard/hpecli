package password

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// ReadFromConsole write the prompt to the console and then
// will read the password from the user console without
// echoing the password to the console
func ReadFromConsole(prompt string) (string, error) {
	fmt.Fprint(os.Stdout, prompt)

	//nolint:unconvert  // false positive lint err
	p, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	// rather than force everyone to do a linefeed after reading
	// we'll just do it here
	fmt.Fprint(os.Stdout, "\n")

	return string(p), nil
}

// ReadFromStdIn reads from os.Stdin and trims any  CR & LF characters
// as well as any spaces around the input.  This method allows passwords
// to be passed to the CLI with commands like:
//   <Windows>
//     echo mypassword | hpecli ilo login --host 1.1.1.1 -u admin --password-stdin
func ReadFromStdIn() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	t, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading from stdin: %w", err)
	}

	t = strings.TrimSpace(t)

	return t, nil
}

func Read(password *string, promptFor bool, prompt string) (err error) {
	if *password != "" && promptFor {
		return errors.New("--password and --password-stdin are mutually exclusive")
	}

	if *password != "" {
		// password specified.. leave it and don't get another one
		return
	}

	// asked to read from stdin
	if promptFor {
		*password, err = ReadFromStdIn()
		return
	}

	// prompt the user for the password on the console
	*password, err = ReadFromConsole(prompt)

	return
}
