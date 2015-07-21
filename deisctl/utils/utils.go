// Package utils contains commonly useful functions from Deisctl
package utils

import (
	"os"
	"strings"
)

// ResolvePath returns the path with a tilde (~) and $HOME replaced by the actual home directory
func ResolvePath(path string) string {
	path = strings.Replace(path, "~", os.Getenv("HOME"), -1)
	// Using $HOME seems to work just fine with `deisctl config`, but not `deisctl refresh-units`
	path = strings.Replace(path, "$HOME", os.Getenv("HOME"), -1)
	return path
}
