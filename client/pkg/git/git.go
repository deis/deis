package git

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

// CreateRemote adds a git remote in the current directory.
func CreateRemote(host, remote, appID string) error {
	cmd := exec.Command("git", "remote", "add", remote, RemoteURL(host, appID))
	stderr, err := cmd.StderrPipe()

	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	output, _ := ioutil.ReadAll(stderr)
	fmt.Print(string(output))

	if err := cmd.Wait(); err != nil {
		return err
	}

	fmt.Printf("Git remote %s added\n", remote)

	return nil
}

// DeleteRemote removes a git remote in the current directory.
func DeleteRemote(appID string) error {
	name, err := remoteNameFromAppID(appID)

	if err != nil {
		return err
	}

	if _, err = exec.Command("git", "remote", "remove", name).Output(); err != nil {
		return err
	}

	fmt.Printf("Git remote %s removed\n", name)

	return nil
}

func remoteNameFromAppID(appID string) (string, error) {
	out, err := exec.Command("git", "remote", "-v").Output()

	if err != nil {
		return "", err
	}

	cmd := string(out)

	for _, line := range strings.Split(cmd, "\n") {
		if strings.Contains(line, appID) {
			return strings.Split(line, "\t")[0], nil
		}
	}

	return "", errors.New("Could not find remote matching app in 'git remote -v'")
}

// DetectAppName detects if there is deis remote in git.
func DetectAppName(host string) (string, error) {
	remote, err := findRemote(host)

	// Don't return an error if remote can't be found, return directory name instead.
	if err != nil {
		dir, err := os.Getwd()
		return strings.ToLower(path.Base(dir)), err
	}

	ss := strings.Split(remote, "/")
	return strings.Split(ss[len(ss)-1], ".")[0], nil
}

func findRemote(host string) (string, error) {
	out, err := exec.Command("git", "remote", "-v").Output()

	if err != nil {
		return "", err
	}

	cmd := string(out)

	// Strip off any trailing :port number after the host name.
	host = strings.Split(host, ":")[0]

	for _, line := range strings.Split(cmd, "\n") {
		for _, remote := range strings.Split(line, " ") {
			if strings.Contains(remote, host) {
				return strings.Split(remote, "\t")[1], nil
			}
		}
	}

	return "", errors.New("Could not find deis remote in 'git remote -v'")
}

// RemoteURL returns the git URL of app.
func RemoteURL(host, appID string) string {
	// Strip off any trailing :port number after the host name.
	host = strings.Split(host, ":")[0]
	return fmt.Sprintf("ssh://git@%s:2222/%s.git", host, appID)
}
