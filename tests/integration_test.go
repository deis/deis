// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/utils"
)

var (
	gitCloneCmd  = "if [ ! -d {{.ExampleApp}} ] ; then git clone https://github.com/deis/{{.ExampleApp}}.git ; fi"
	gitRemoveCmd = "git remote remove deis"
	gitPushCmd   = "git push deis master"
	gitAddCmd    = "git add ."
	gitCommitCmd = "git commit -m fake"
)

func TestGlobal(t *testing.T) {
	params := utils.GetGlobalConfig()
	cookieTest(t, params)
	utils.Execute(t, authRegisterCmd, params, false, "")
	utils.Execute(t, keysAddCmd, params, false, "")
	utils.Execute(t, clustersCreateCmd, params, false, "")
}

func cookieTest(t *testing.T, params *utils.DeisTestConfig) {
	// Regression test for https://github.com/deis/deis/pull/1136
	// Ensure that cookies are cleared on auth:register and auth:cancel
	utils.Execute(t, authRegisterCmd, params, false, "")
	cmd := "cat ~/.deis/cookies.txt"
	utils.CheckList(t, cmd, params, "csrftoken", false)
	utils.CheckList(t, cmd, params, "sessionid", false)
	utils.AuthCancel(t, params)
	utils.CheckList(t, cmd, params, "csrftoken", true)
	utils.CheckList(t, cmd, params, "sessionid", true)
}
