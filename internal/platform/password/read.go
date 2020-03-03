package password

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// Password write the prompt to the console and then
// will read the password from the user console without
// echoing the password to the console
func ReadFromConsole(prompt string) (string, error) {
	fmt.Fprint(os.Stdout, prompt)

	p, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	return string(p), nil
}

// ReadFromEnvOrConsole first looks in the environment for the
// "envar" parameter.  If set, then it returns the value
// as the password.  If then env variable is not found, it will
// then print the prompt to the console and read the password
// from the console.
func ReadFromEnvOrConsole(prompt, envar string) (string, error) {
	p := os.Getenv(envar)
	if p != "" {
		return p, nil
	}

	return ReadFromConsole(prompt)
}
