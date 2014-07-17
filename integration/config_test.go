package verbose

import (
	_ "fmt"
	"github.com/deis/deis/tests/integration-utils"
	"github.com/deis/deis/tests/utils"
	"testing"
)

func configSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	cfg.ExampleApp = itutils.GetRandomApp()
	cfg.AppName = "configsample"
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

func configlistTest(t *testing.T, params *itutils.DeisTestConfig, notflag bool) {
	cmd := itutils.GetCommand("config", "list")
	itutils.CheckList(t, params, cmd, "jaf", notflag)

}

func configSetTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("config", "set")
	itutils.Execute(t, cmd, params, false, "")
	itutils.CheckList(t, params, itutils.GetCommand("apps", "info"), "(v3)", false)
}

func configUnsetTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("config", "unset")
	itutils.Execute(t, cmd, params, false, "")
	itutils.CheckList(t, params, itutils.GetCommand("apps", "info"), "(v4)", false)
}

func appsOpenTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Curl(t, params)
}

func TestConfig(t *testing.T) {
	params := configSetup(t)
	configSetTest(t, params)
	configlistTest(t, params, false)
	appsOpenTest(t, params)
	configUnsetTest(t, params)
	configlistTest(t, params, true)
	itutils.AppsDestroyTest(t, params)
}
