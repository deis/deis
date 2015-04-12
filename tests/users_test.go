// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/utils"
)

var (
	usersListCmd = "users:list"
)

func TestUsers(t *testing.T) {
	params := utils.GetGlobalConfig()
	user := utils.GetGlobalConfig()
	user.UserName, user.Password = "user-list-test", "test"
	user.AppName = params.AppName
	utils.Execute(t, authRegisterCmd, user, false, "")
	usersListTest(t, params, user)
}

func usersListTest(t *testing.T, params *utils.DeisTestConfig, user *utils.DeisTestConfig) {
	utils.Execute(t, authLoginCmd, user, false, "")
	utils.Execute(t, usersListCmd, user, true, "403 FORBIDDEN")
	utils.Execute(t, authLoginCmd, params, false, "")
	utils.CheckList(t, usersListCmd, params, "user-list-test", false)
}
