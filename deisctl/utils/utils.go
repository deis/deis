// Package utils contains commonly useful functions from Deisctl
package utils

import (
	"fmt"
	"os"
	"strings"
)

// DeisIfy returns a pretty-printed deis logo along with the corresponding message
func DeisIfy(message string) string {
	circle := "\033[31m●"
	square := "\033[32m■"
	triangle := "\033[34m▴"
	reset := "\033[0m"
	title := reset + message

	return fmt.Sprintf("%s %s %s\n%s %s %s %s\n%s %s %s%s\n",
		circle, triangle, square,
		square, circle, triangle, title,
		triangle, square, circle, reset)
}

// ResolvePath returns the path with a tilde (~) and $HOME replaced by the actual home directory
func ResolvePath(path string) string {
	path = strings.Replace(path, "~", os.Getenv("HOME"), -1)
	// Using $HOME seems to work just fine with `deisctl config`, but not `deisctl refresh-units`
	path = strings.Replace(path, "$HOME", os.Getenv("HOME"), -1)
	return path
}
