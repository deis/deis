// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/utils"
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
func keysSetup(t *testing.T) *utils.DeisTestConfig {
	cfg := utils.GetGlobalConfig()
	utils.Execute(t, authLoginCmd, cfg, false, "")
	return cfg
}

func keysAddTest(t *testing.T, params *utils.DeisTestConfig) {
	params.AuthKey = "deiskey"
	utils.Execute(t, keysAddCmd, params, false, "")
	utils.Execute(t, keysAddCmd, params, true,
		"This field must be unique")
}

func keysListTest(t *testing.T, params *utils.DeisTestConfig, notflag bool) {
	utils.CheckList(t, keysListCmd, params, params.AuthKey, notflag)
}

func keysRemoveTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.Execute(t, keysRemoveCmd, params, false, "")
}
