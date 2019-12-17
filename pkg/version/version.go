// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package version

import "fmt"

// Get returns the short version. just the version (e.g. v0.1)
func Get() string {
	return version
}

// GetFull returns the long version. (e.g. v0.2:6683f37:2019-11-23)
func GetFull() string {
	return fmt.Sprintf("%s:%s:%s", version, gitCommit, builtAt)
}
