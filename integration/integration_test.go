package verbose

import (
	"github.com/deis/deis/tests/integration-utils"
	"testing"
)

func authcancelTest(t *testing.T, params *itutils.DeisTestConfig) {
	var cmd string
	cmd = itutils.GetCommand("auth", "register")
	itutils.Execute(t, cmd, params, false, "")
	itutils.AuthCancel(t, params)
}

func TestGlobal(t *testing.T) {
	params := itutils.GetGlobalConfig()
	var cmd string
	authcancelTest(t, params)
	cmd = itutils.GetCommand("auth", "register")
	itutils.Execute(t, cmd, params, false, "")
	cmd = itutils.GetCommand("keys", "add")
	itutils.Execute(t, cmd, params, false, "")
	cmd = itutils.GetCommand("clusters", "create")
	itutils.Execute(t, cmd, params, false, "")
}
