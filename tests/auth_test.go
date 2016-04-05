// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/utils"
)

var (
	authLoginCmd         = "auth:login http://deis.{{.Domain}} --username={{.UserName}} --password={{.Password}}"
	authLogoutCmd        = "auth:logout"
	authRegisterCmd      = "auth:register http://deis.{{.Domain}} --username={{.UserName}} --password={{.Password}} --email={{.Email}}"
	authCancelCmd        = "auth:cancel --username={{.UserName}} --password={{.Password}} --yes"
	authCancelAdminCmd   = "auth:cancel --username={{.UserName}} --yes"
	authRegenerateCmd    = "auth:regenerate"
	authRegenerateUsrCmd = "auth:regenerate -u {{.UserName}}"
	authRegenerateAllCmd = "auth:regenerate --all"
	checkTokenCmd        = "apps:list"
	authPasswdCmd        = "auth:passwd --username={{.UserName}} --password={{.Password}} --new-password={{.NewPassword}}"
	authWhoamiCmd        = "auth:whoami"
)

func TestAuth(t *testing.T) {
	params := authSetup(t)
	authRegisterTest(t, params)
	authLogoutTest(t, params)
	authRegenerateTest(t)
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
	utils.Execute(t, authCancelCmd, params, false, "Account cancelled")
	user := utils.GetGlobalConfig()

	// Admins can delete other users.
	user.UserName, user.Password = "cancel-test", "test"
	utils.Execute(t, authRegisterCmd, user, false, "")
	admin := utils.GetGlobalConfig()
	utils.Execute(t, authLoginCmd, admin, false, "")
	utils.Execute(t, authCancelAdminCmd, user, false, "Account cancelled")
	// Make sure the user's config was not purged after auth:cancel --username
	utils.Execute(t, authWhoamiCmd, admin, false, "You are "+admin.UserName)
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
	params.NewPassword = password
	cmd := authPasswdCmd
	utils.Execute(t, cmd, params, false, "")
	params.Password = "wrong-password"
	utils.Execute(t, cmd, params, true, "Password change failed")

	cmd = authLoginCmd
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

func authRegenerateTest(t *testing.T) {
	params := utils.GetGlobalConfig()
	regenCmd := authRegenerateUsrCmd
	loginCmd := authLoginCmd

	utils.Execute(t, loginCmd, params, false, "")
	utils.Execute(t, authRegenerateCmd, params, false, "")
	utils.Execute(t, loginCmd, params, false, "")
	utils.Execute(t, regenCmd, params, false, "")
	utils.Execute(t, checkTokenCmd, params, true, "401 UNAUTHORIZED")
	utils.Execute(t, loginCmd, params, false, "")
	utils.Execute(t, authRegenerateAllCmd, params, false, "")
	utils.Execute(t, checkTokenCmd, params, true, "401 UNAUTHORIZED")
}
