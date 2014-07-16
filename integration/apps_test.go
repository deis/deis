package verbose

import (
	"github.com/deis/deis/tests/integration-utils"
	"github.com/deis/deis/tests/utils"
	"testing"
)

func appsSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	cfg.ExampleApp = itutils.GetRandomApp()
	cmd := itutils.GetCommand("auth", "login")
	itutils.Execute(t, cmd, cfg, false, "")
	cmd = itutils.GetCommand("git", "clone")
	itutils.Execute(t, cmd, cfg, false, "")
	return cfg
}

func appsCreateTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("apps", "create")
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
	itutils.Execute(t, cmd, params, false, "")
	itutils.Execute(t, cmd, params, true, "Deis remote already exists")

	if err := utils.Chdir(".."); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
}

func appsRunTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("apps", "run")
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
	itutils.Execute(t, cmd, params, false, "")

	if err := utils.Chdir(".."); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
	itutils.Execute(t, cmd, params, true, "Could not find deis remote in `git remote -v`")
}

func appsDestroyTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("apps", "destroy")
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
	itutils.Execute(t, cmd, params, false, "")
	itutils.Execute(t, cmd, params, true, "400 BAD REQUEST")
	if err := utils.Chdir(".."); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
	if err := utils.Rmdir(params.ExampleApp); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
}

func appsListTest(t *testing.T, params *itutils.DeisTestConfig, notflag bool) {
	cmd := itutils.GetCommand("apps", "list")
	itutils.CheckList(t, params, cmd, params.AppName, notflag)
}

func appsLogsTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("apps", "logs")
	cmd1 := itutils.GetCommand("git", "push")
	itutils.Execute(t, cmd, params, true, "204 NO CONTENT")
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
	itutils.Execute(t, cmd1, params, false, "")
	itutils.Execute(t, cmd, params, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
}

func appsInfoTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("apps", "info")
	itutils.Execute(t, cmd, params, false, "")
}

func appsOpenTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Curl(t, "http://"+params.AppName+"."+params.HostName, params.ExampleApp)
}

func TestApps(t *testing.T) {
	params := appsSetup(t)
	appsCreateTest(t, params)
	appsListTest(t, params, false)
	appsLogsTest(t, params)
	appsInfoTest(t, params)
	appsRunTest(t, params)
	appsOpenTest(t, params)
	appsDestroyTest(t, params)
	appsListTest(t, params, true)

}
