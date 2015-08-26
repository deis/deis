package cmd

import (
	"github.com/deis/deis/client/pkg/git"
)

// GitRemote creates a git remote for a deis app.
func GitRemote(appID, remote string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	return git.CreateRemote(c.ControllerURL.Host, remote, appID)
}
