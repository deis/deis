// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/integration-utils"
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
	user := itutils.GetGlobalConfig()
	user.UserName, user.Password = "test1", "test1"
	user.AppName = params.AppName
	itutils.Execute(t, authRegisterCmd, user, false, "")
	permsCreateAppTest(t, params, user)
	permsDeleteAppTest(t, params, user)
	permsCreateAdminTest(t, params)
	permsDeleteAdminTest(t, params)
	itutils.AppsDestroyTest(t, params)
}

func permsSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	cfg.AppName = "permssample"
	itutils.Execute(t, authLoginCmd, cfg, false, "")
	itutils.Execute(t, gitCloneCmd, cfg, false, "")
	if err := utils.Chdir(cfg.ExampleApp); err != nil {
		t.Fatal(err)
	}
	itutils.Execute(t, appsCreateCmd, cfg, false, "")
	itutils.Execute(t, gitPushCmd, cfg, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	return cfg
}

func permsCreateAdminTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Execute(t, permsCreateAdminCmd, params, false, "")
	itutils.CheckList(t, permsListAdminCmd, params, "test1", false)
}

func permsCreateAppTest(t *testing.T, params, user *itutils.DeisTestConfig) {
	itutils.Execute(t, authLoginCmd, user, false, "")
	itutils.Execute(t, permsCreateAppCmd, user, true, "403 FORBIDDEN")
	itutils.Execute(t, authLoginCmd, params, false, "")
	itutils.Execute(t, permsCreateAppCmd, params, false, "")
	itutils.CheckList(t, permsListAppCmd, params, "test1", false)
}

func permsDeleteAdminTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Execute(t, permsDeleteAdminCmd, params, false, "")
	itutils.CheckList(t, permsListAdminCmd, params, "test1", true)
}

func permsDeleteAppTest(t *testing.T, params, user *itutils.DeisTestConfig) {
	itutils.Execute(t, authLoginCmd, user, false, "")
	itutils.Execute(t, permsDeleteAppCmd, user, true, "403 FORBIDDEN")
	itutils.Execute(t, authLoginCmd, params, false, "")
	itutils.Execute(t, permsDeleteAppCmd, params, false, "")
	itutils.CheckList(t, permsListAppCmd, params, "test1", true)
}
