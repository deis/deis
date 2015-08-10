package webbrowser

import (
	"os/exec"
)

// Webbrowser opens a URL with the default browser.
func Webbrowser(u string) (err error) {
	_, err = exec.Command("cmd", "/c", "start", u).Output()
	return
}
