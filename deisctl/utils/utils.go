// Package utils contains commonly useful functions from Deisctl

package utils

import "fmt"

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
