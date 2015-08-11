// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/utils"
)

var (
	permsListAppCmd     = "perms:list --app={{.AppName}}"
	permsListAdminCmd   = "perms:list --admin"
	permsCreateAppCmd   = "perms:create {{.AppUser}} --app={{.AppName}}"
	permsCreateAdminCmd = "perms:create {{.AppUser}} --admin"
	permsDeleteAppCmd   = "perms:delete {{.AppUser}} --app={{.AppName}}"
	permsDeleteAdminCmd = "perms:delete {{.AppUser}} --admin"
)

func TestPerms(t *testing.T) {
	params := permsSetup(t)
	user := utils.GetGlobalConfig()
	user.UserName, user.Password = "test1", "test1"
	user.AppName = params.AppName
	utils.Execute(t, authRegisterCmd, user, false, "")
	permsCreateAppTest(t, params, user)
	permsDeleteAppTest(t, params, user)
	permsCreateAdminTest(t, params)
	permsDeleteAdminTest(t, params)
	utils.AppsDestroyTest(t, params)
}

func permsSetup(t *testing.T) *utils.DeisTestConfig {
	cfg := utils.GetGlobalConfig()
	cfg.AppName = "permssample"
	utils.Execute(t, authLoginCmd, cfg, false, "")
	utils.Execute(t, gitCloneCmd, cfg, false, "")
	if err := utils.Chdir(cfg.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, appsCreateCmd, cfg, false, "")
	utils.Execute(t, gitPushCmd, cfg, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	return cfg
}

func permsCreateAdminTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.Execute(t, permsCreateAdminCmd, params, false, "")
	utils.CheckList(t, permsListAdminCmd, params, "test1", false)
}

func permsCreateAppTest(t *testing.T, params, user *utils.DeisTestConfig) {
	utils.Execute(t, authLoginCmd, user, false, "")
	utils.Execute(t, permsCreateAppCmd, user, true, "403 FORBIDDEN")
	utils.Execute(t, authLoginCmd, params, false, "")
	utils.Execute(t, permsCreateAppCmd, params, false, "")
	utils.CheckList(t, permsListAppCmd, params, "test1", false)
}

func permsDeleteAdminTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.Execute(t, permsDeleteAdminCmd, params, false, "")
	utils.CheckList(t, permsListAdminCmd, params, "test1", true)
}

func permsDeleteAppTest(t *testing.T, params, user *utils.DeisTestConfig) {
	utils.Execute(t, authLoginCmd, user, false, "")
	utils.Execute(t, permsDeleteAppCmd, user, false, "")
	utils.Execute(t, authLoginCmd, params, false, "")
	utils.CheckList(t, permsListAppCmd, params, "test1", true)
}
