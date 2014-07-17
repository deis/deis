package verbose

import (
	"github.com/deis/deis/tests/integration-utils"
	"testing"
)

//Tests #1136 // Tests #1239
func cookieTest(t *testing.T, params *itutils.DeisTestConfig) {
	var cmd string
	cmd = itutils.GetCommand("auth", "register")
	itutils.Execute(t, cmd, params, false, "")
	itutils.CheckList(t, params, "cat ~/.deis/cookies.txt", "csrftoken", false)
	itutils.CheckList(t, params, "cat ~/.deis/cookies.txt", "sessionid", false)
	itutils.AuthCancel(t, params)
	itutils.CheckList(t, params, "cat ~/.deis/cookies.txt", "csrftoken", true)
	itutils.CheckList(t, params, "cat ~/.deis/cookies.txt", "sessionid", true)
}

func TestGlobal(t *testing.T) {
	params := itutils.GetGlobalConfig()
	var cmd string
	cookieTest(t, params)
	cmd = itutils.GetCommand("auth", "register")
	itutils.Execute(t, cmd, params, false, "")
	cmd = itutils.GetCommand("keys", "add")
	itutils.Execute(t, cmd, params, false, "")
	cmd = itutils.GetCommand("clusters", "create")
	itutils.Execute(t, cmd, params, false, "")
}
