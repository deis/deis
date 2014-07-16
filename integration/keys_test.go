package verbose

import (
	_ "fmt"
	"github.com/deis/deis/tests/integration-utils"
	"testing"
)

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
	itutils.Execute(t, cmd, params, true, "Uploading deiskey to Deis...400 BAD REQUEST")
}

func keysListTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("keys", "list")
	itutils.CheckList(t, params, cmd, params.ClusterName, notflag)
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
	keysListTest(t, params, false)
}
