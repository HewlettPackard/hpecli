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

	//nolint:unconvert  // false positive lint err
	p, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	return string(p), nil
}
