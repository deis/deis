// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/integration-utils"
)

// Requires a ~/.ssh/deis-testkey to be set up:
// $ ssh-keygen -q -t rsa -f ~/.ssh/deiskey -N '' -C deiskey
func keysSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	cmd := itutils.GetCommand("auth", "login")
	itutils.Execute(t, cmd, cfg, false, "")
	return cfg
}

func keysAddTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("keys", "add")
	params.AuthKey = "deiskey"
	itutils.Execute(t, cmd, params, false, "")
	itutils.Execute(t, cmd, params, true, "SSH Key with this Public already exists")
}

func keysListTest(t *testing.T, params *itutils.DeisTestConfig, notflag bool) {
	cmd := itutils.GetCommand("keys", "list")
	itutils.CheckList(t, params, cmd, params.AuthKey, notflag)
}

func keysRemoveTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("keys", "remove")
	itutils.Execute(t, cmd, params, false, "")
}

func TestKeys(t *testing.T) {
	params := keysSetup(t)
	keysAddTest(t, params)
	keysListTest(t, params, false)
	keysRemoveTest(t, params)
	keysListTest(t, params, true)
}
