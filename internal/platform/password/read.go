package password

import (
	"fmt"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// Password will read the password from the user console without echoing
// the password to teh console
func Read(prompt string) (string, error) {
	fmt.Print(prompt)
	p, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	return string(p), nil
}
