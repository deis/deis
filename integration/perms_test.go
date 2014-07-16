package verbose

import (
	_ "fmt"
	"github.com/deis/deis/tests/integration-utils"
	"github.com/deis/deis/tests/utils"
	"testing"
)

func permsSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	cfg.ExampleApp = itutils.GetRandomApp()
	cfg.AppName = "permssample"
	cmd := itutils.GetCommand("auth", "login")
	itutils.Execute(t, cmd, cfg, false, "")
	cmd = itutils.GetCommand("git", "clone")
	itutils.Execute(t, cmd, cfg, false, "")
	cmd = itutils.GetCommand("apps", "create")
	cmd1 := itutils.GetCommand("git", "push")
	if err := utils.Chdir(cfg.ExampleApp); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}

	itutils.Execute(t, cmd, cfg, false, "")
	itutils.Execute(t, cmd1, cfg, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
	return cfg
}

func permsCreateAppTest(t *testing.T, params, user *itutils.DeisTestConfig) {
	var cmd string
	cmd = itutils.GetCommand("auth", "login")
	itutils.Execute(t, cmd, user, false, "")
	cmd = itutils.GetCommand("perms", "create-app")
	itutils.Execute(t, cmd, user, true, "403 FORBIDDEN")
	cmd = itutils.GetCommand("auth", "login")
	itutils.Execute(t, cmd, params, false, "")
	cmd = itutils.GetCommand("perms", "create-app")
	itutils.Execute(t, cmd, params, false, "")
	cmd = itutils.GetCommand("perms", "list-app")
	itutils.CheckList(t, params, cmd, "test1", false)
}

func permsDeleteAppTest(t *testing.T, params, user *itutils.DeisTestConfig) {
	var cmd string
	cmd = itutils.GetCommand("auth", "login")
	itutils.Execute(t, cmd, user, false, "")
	cmd = itutils.GetCommand("perms", "delete-app")
	itutils.Execute(t, cmd, user, true, "403 FORBIDDEN")
	cmd = itutils.GetCommand("auth", "login")
	itutils.Execute(t, cmd, params, false, "")
	cmd = itutils.GetCommand("perms", "delete-app")
	itutils.Execute(t, cmd, params, false, "")
	cmd = itutils.GetCommand("perms", "list-app")
	itutils.CheckList(t, params, cmd, "test1", true)
}

func permsCreateAdminTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("perms", "create-admin")
	itutils.Execute(t, cmd, params, false, "")
	cmd = itutils.GetCommand("perms", "list-admin")
	itutils.CheckList(t, params, cmd, "test1", false)

}

func permsDeleteAdminTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("perms", "delete-admin")
	itutils.Execute(t, cmd, params, false, "")
	cmd = itutils.GetCommand("perms", "list-admin")
	itutils.CheckList(t, params, cmd, "test1", true)
}

func TestBuilds(t *testing.T) {
	params := permsSetup(t)
	user := itutils.GetGlobalConfig()
	user.UserName, user.Password = "test1", "test1"
	user.AppName = params.AppName
	cmd := itutils.GetCommand("auth", "register")
	itutils.Execute(t, cmd, user, false, "")
	permsCreateAppTest(t, params, user)
	permsDeleteAppTest(t, params, user)
	permsCreateAdminTest(t, params)
	permsDeleteAdminTest(t, params)
	itutils.AppsDestroyTest(t, params)
}
