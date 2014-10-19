// +build integration

package tests

import (
	"os"
	"os/user"
	"testing"

	"github.com/deis/deis/tests/utils"
)

var (
	gitCloneCmd  = "if [ ! -d {{.ExampleApp}} ] ; then git clone https://github.com/deis/{{.ExampleApp}}.git ; fi"
	gitRemoveCmd = "git remote remove deis"
	gitPushCmd   = "git push deis master"
)

func TestGlobal(t *testing.T) {
	params := utils.GetGlobalConfig()
	cookieTest(t, params)
	utils.Execute(t, authRegisterCmd, params, false, "")
	utils.Execute(t, keysAddCmd, params, false, "")
}

func cookieTest(t *testing.T, params *utils.DeisTestConfig) {
	// Regression test for https://github.com/deis/deis/pull/1136
	// Ensure that cookies are cleared on auth:register and auth:cancel
	user, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	cookieJar := user.HomeDir + "/.deis/cookies.txt"
	utils.Execute(t, authRegisterCmd, params, false, "")
	cmd := "cat " + cookieJar
	utils.CheckList(t, cmd, params, "csrftoken", false)
	utils.CheckList(t, cmd, params, "sessionid", false)
	info, err := os.Stat(cookieJar)
	if err != nil {
		t.Fatal(err)
	}
	mode := info.Mode().String()
	expected := "-rw-------"
	if mode != expected {
		t.Fatalf("%s has wrong mode:\n   current: %s\n  expected: %s",
			cookieJar, mode, expected)
	}
	utils.AuthCancel(t, params)
	utils.CheckList(t, cmd, params, "csrftoken", true)
	utils.CheckList(t, cmd, params, "sessionid", true)
}
