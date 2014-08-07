// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/integration-utils"
)

var (
	gitCloneCmd  = "if [ ! -d {{.ExampleApp}} ] ; then git clone https://github.com/deis/{{.ExampleApp}}.git ; fi"
	gitRemoveCmd = "git remote remove deis"
	gitPushCmd   = "git push deis master"
	gitAddCmd    = "git add ."
	gitCommitCmd = "git commit -m fake"
)

func TestGlobal(t *testing.T) {
	params := itutils.GetGlobalConfig()
	cookieTest(t, params)
	itutils.Execute(t, authRegisterCmd, params, false, "")
	itutils.Execute(t, keysAddCmd, params, false, "")
	itutils.Execute(t, clustersCreateCmd, params, false, "")
}

func cookieTest(t *testing.T, params *itutils.DeisTestConfig) {
	// Regression test for https://github.com/deis/deis/pull/1136
	// Ensure that cookies are cleared on auth:register and auth:cancel
	itutils.Execute(t, authRegisterCmd, params, false, "")
	cmd := "cat ~/.deis/cookies.txt"
	itutils.CheckList(t, cmd, params, "csrftoken", false)
	itutils.CheckList(t, cmd, params, "sessionid", false)
	itutils.AuthCancel(t, params)
	itutils.CheckList(t, cmd, params, "csrftoken", true)
	itutils.CheckList(t, cmd, params, "sessionid", true)
}
