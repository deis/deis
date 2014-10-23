// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/utils"
)

var (
	authLoginCmd    = "auth:login http://deis.{{.Domain}} --username={{.UserName}} --password={{.Password}}"
	authLogoutCmd   = "auth:logout"
	authRegisterCmd = "auth:register http://deis.{{.Domain}} --username={{.UserName}} --password={{.Password}} --email={{.Email}}"
)

func TestAuth(t *testing.T) {
	params := authSetup(t)
	authRegisterTest(t, params)
	authLogoutTest(t, params)
	authLoginTest(t, params)
	authWhoamiTest(t, params)
	authPasswdTest(t, params)
	authCancel(t, params)
}

func authSetup(t *testing.T) *utils.DeisTestConfig {
	user := utils.GetGlobalConfig()
	user.UserName, user.Password = utils.NewID(), utils.NewID()
	return user
}

func authCancel(t *testing.T, params *utils.DeisTestConfig) {
	utils.AuthCancel(t, params)
}

func authLoginTest(t *testing.T, params *utils.DeisTestConfig) {
	cmd := authLoginCmd
	utils.Execute(t, cmd, params, false, "")
	params = authSetup(t)
	utils.Execute(t, cmd, params, true, "400 BAD REQUEST")
}

func authLogoutTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.Execute(t, authLogoutCmd, params, false, "")
}

func authPasswdTest(t *testing.T, params *utils.DeisTestConfig) {
	password := "aNewPassword"
	utils.AuthPasswd(t, params, password)
	cmd := authLoginCmd
	utils.Execute(t, cmd, params, true, "400 BAD REQUEST")
	params.Password = password
	utils.Execute(t, cmd, params, false, "")
}

func authRegisterTest(t *testing.T, params *utils.DeisTestConfig) {
	cmd := authRegisterCmd
	utils.Execute(t, cmd, params, false, "")
	utils.Execute(t, cmd, params, true, "Registration failed")
}

func authWhoamiTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.Execute(t, "auth:whoami", params, true, params.UserName)
}
