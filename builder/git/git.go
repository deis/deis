package git

// This file just contains the Git-specific portions of sshd.

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/cookoo/log"
	"golang.org/x/crypto/ssh"
)

// PrereceiveHookTmpl is a pre-recieve hook.
//
// This is overridable. The following template variables are passed into it:
//
// 	.GitHome: the path to Git's home directory.
var PrereceiveHookTpl = `#!/bin/bash
strip_remote_prefix() {
    stdbuf -i0 -o0 -e0 sed "s/^/"$'\e[1G'"/"
}

while read oldrev newrev refname
do
  LOCKFILE="/tmp/$RECEIVE_REPO.lock"
  if ( set -o noclobber; echo "$$" > "$LOCKFILE" ) 2> /dev/null; then
	trap 'rm -f "$LOCKFILE"; exit 1' INT TERM EXIT

	# check for authorization on this repo
	{{.GitHome}}/receiver "$RECEIVE_REPO" "$newrev" "$RECEIVE_USER" "$RECEIVE_FINGERPRINT"
	rc=$?
	if [[ $rc != 0 ]] ; then
	  echo "      ERROR: failed on rev $newrev - push denied"
	  exit $rc
	fi
	# builder assumes that we are running this script from $GITHOME
	cd {{.GitHome}}
	# if we're processing a receive-pack on an existing repo, run a build
	if [[ $SSH_ORIGINAL_COMMAND == git-receive-pack* ]]; then
		{{.GitHome}}/builder "$RECEIVE_USER" "$RECEIVE_REPO" "$newrev" 2>&1 | strip_remote_prefix
	fi

	rm -f "$LOCKFILE"
	trap - INT TERM EXIT
  else
	echo "Another git push is ongoing. Aborting..."
	exit 1
  fi
done
`

// Receive receives a Git repo.
// This will only work for git-receive-pack.
//
// Params:
// 	- operation (string): e.g. git-receive-pack
// 	- repoName (string): The repository name, in the form '/REPO.git'.
// 	- channel (ssh.Channel): The channel.
// 	- request (*ssh.Request): The channel.
// 	- gitHome (string): Defaults to /home/git.
// 	- fingerprint (string): The fingerprint of the user's SSH key.
// 	- user (string): The name of the Deis user.
//
// Returns:
// 	- nothing
func Receive(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	if ok, z := p.Requires("channel", "request", "fingerprint", "permissions"); !ok {
		return nil, fmt.Errorf("Missing requirements %q", z)
	}
	repoName := p.Get("repoName", "").(string)
	operation := p.Get("operation", "").(string)
	channel := p.Get("channel", nil).(ssh.Channel)
	//req := p.Get("request", nil).(*ssh.Request)
	gitHome := p.Get("gitHome", "/home/git").(string)
	fingerprint := p.Get("fingerprint", nil).(string)
	user := p.Get("user", "").(string)

	repo, err := cleanRepoName(repoName)
	if err != nil {
		log.Warnf(c, "Illegal repo name: %s.", err)
		channel.Stderr().Write([]byte("No repo given"))
		return nil, err
	}
	repo += ".git"

	if _, err := createRepo(c, filepath.Join(gitHome, repo), gitHome); err != nil {
		log.Infof(c, "Did not create new repo: %s", err)
	}
	cmd := exec.Command("git-shell", "-c", fmt.Sprintf("%s '%s'", operation, repo))
	//cmd := execAs("git", "/usr/bin/git-shell", "-c", fmt.Sprintf("%s '%s'", operation, repo))
	log.Infof(c, strings.Join(cmd.Args, " "))

	var errbuff bytes.Buffer
	// errwriter := io.MultiWriter(channel.Stderr(), &errbuff)

	cmd.Dir = gitHome
	cmd.Env = []string{
		fmt.Sprintf("RECEIVE_USER=%s", user),
		fmt.Sprintf("RECEIVE_REPO=%s", repo),
		fmt.Sprintf("RECEIVE_FINGERPRINT=%s", fingerprint),
		fmt.Sprintf("SSH_ORIGINAL_COMMAND=%s '%s'", operation, repo),
		fmt.Sprintf("SSH_CONNECTION=%s", c.Get("SSH_CONNECTION", "0 0 0 0").(string)),
	}
	cmd.Env = append(cmd.Env, os.Environ()...)

	done := plumbCommand(cmd, channel, &errbuff)

	//cmd.Stdin = channel
	if err := cmd.Start(); err != nil {
		log.Warnf(c, "Failed git receive immediately: %s %s", err, errbuff.Bytes())
		return nil, err
	}
	fmt.Printf("Waiting for git-receive to run.\n")
	done.Wait()
	fmt.Printf("Waiting for deploy.\n")
	if err := cmd.Wait(); err != nil {
		log.Errf(c, "Error on command: %s %s", err, errbuff.Bytes())
		return nil, err
	}
	if errbuff.Len() > 0 {
		log.Warnf(c, "Unreported error: %s", errbuff.Bytes())
	}
	log.Infof(c, "Deploy complete.\n")

	return nil, nil
}

func execAs(user, cmd string, args ...string) *exec.Cmd {
	fullCmd := cmd + " " + strings.Join(args, " ")
	return exec.Command("su", user, "-c", fullCmd)
}

// cleanRepoName cleans a repository name for a git-sh operation.
func cleanRepoName(name string) (string, error) {
	if len(name) == 0 {
		return name, errors.New("Empty repo name.")
	}
	if strings.Contains(name, "..") {
		return "", errors.New("Cannot change directory in file name.")
	}
	name = strings.Replace(name, "'", "", -1)
	return strings.TrimPrefix(strings.TrimSuffix(name, ".git"), "/"), nil
}

// plumbCommand connects the exec in/output and the channel in/output.
//
// The sidechannel is for sending errors to logs.
func plumbCommand(cmd *exec.Cmd, channel ssh.Channel, sidechannel io.Writer) *sync.WaitGroup {
	var wg sync.WaitGroup
	inpipe, _ := cmd.StdinPipe()
	go func() {
		io.Copy(inpipe, channel)
		inpipe.Close()
	}()

	cmd.Stdout = channel
	cmd.Stderr = channel.Stderr()

	return &wg
}

var createLock sync.Mutex

// createRepo creates a new Git repo if it is not present already.
//
// Largely inspired by gitreceived from Flynn.
//
// Returns a bool indicating whether a project was created (true) or already
// existed (false).
func createRepo(c cookoo.Context, repoPath, gitHome string) (bool, error) {
	createLock.Lock()
	defer createLock.Unlock()

	if fi, err := os.Stat(repoPath); err == nil && fi.IsDir() {
		// Nothing to do.
		log.Infof(c, "Directory %s already exists.", repoPath)
		return false, nil
	} else if os.IsNotExist(err) {

		log.Infof(c, "Creating new directory at %s", repoPath)
		// Create directory
		if err := os.MkdirAll(repoPath, 0755); err != nil {
			log.Warnf(c, "Failed to create repository: %s", err)
			return false, err
		}
		cmd := exec.Command("git", "init", "--bare")
		cmd.Dir = repoPath
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Warnf(c, "git init output: %s", out)
			return false, err
		}

		hook, err := prereceiveHook(map[string]string{"GitHome": gitHome})
		if err != nil {
			return true, err
		}
		ioutil.WriteFile(filepath.Join(repoPath, "hooks", "pre-receive"), hook, 0755)

		return true, nil
	} else if err == nil {
		return false, errors.New("Expected directory, found file.")
	} else {
		return false, err
	}
}

//prereceiveHook templates a pre-receive hook for Git.
func prereceiveHook(vars map[string]string) ([]byte, error) {
	var out bytes.Buffer
	// We parse the template anew each receive in case it has changed.
	t, err := template.New("hooks").Parse(PrereceiveHookTpl)
	if err != nil {
		return []byte{}, err
	}

	err = t.Execute(&out, vars)
	return out.Bytes(), err
}
