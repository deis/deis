// +build linux darwin

package client

import (
	"os"
)

// FindHome returns the HOME directory of the current user
func FindHome() string {
	return os.Getenv("HOME")
}
