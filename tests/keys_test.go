// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/integration-utils"
)

var (
	keysAddCmd    = "keys:add ~/.ssh/{{.AuthKey}}.pub || true"
	keysListCmd   = "keys:list"
	keysRemoveCmd = "keys:remove {{.AuthKey}} || true"
)

func TestKeys(t *testing.T) {
	params := keysSetup(t)
	keysAddTest(t, params)
	keysListTest(t, params, false)
	keysRemoveTest(t, params)
	keysListTest(t, params, true)
}

// Requires a ~/.ssh/deis-testkey to be set up:
// $ ssh-keygen -q -t rsa -f ~/.ssh/deiskey -N '' -C deiskey
func keysSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	itutils.Execute(t, authLoginCmd, cfg, false, "")
	return cfg
}

func keysAddTest(t *testing.T, params *itutils.DeisTestConfig) {
	params.AuthKey = "deiskey"
	itutils.Execute(t, keysAddCmd, params, false, "")
	itutils.Execute(t, keysAddCmd, params, true,
		"SSH Key with this Public already exists")
}

func keysListTest(t *testing.T, params *itutils.DeisTestConfig, notflag bool) {
	itutils.CheckList(t, keysListCmd, params, params.AuthKey, notflag)
}

func keysRemoveTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Execute(t, keysRemoveCmd, params, false, "")
}
