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
)

func TestGlobal(t *testing.T) {
	params := utils.GetGlobalConfig()
	utils.Execute(t, authRegisterCmd, params, false, "")
	utils.Execute(t, keysAddCmd, params, false, "")
}
