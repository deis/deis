package cmd

// GitRemote creates a git remote for a deis app.
func GitRemote(appID, remote string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	return c.CreateRemote(remote, appID)
}
