package webbrowser

import (
	"os/exec"
)

// Webbrowser opens a url with the default browser.
func Webbrowser(u string) (err error) {
	_, err = exec.Command("open", u).Output()
	return
}
