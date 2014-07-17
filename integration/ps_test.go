package verbose

import (
	_ "fmt"
	"github.com/deis/deis/tests/integration-utils"
	"github.com/deis/deis/tests/utils"
	"testing"
)

func psSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	cfg.ExampleApp = itutils.GetRandomApp()
	cfg.AppName = "pssample"
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

func psListTest(t *testing.T, params *itutils.DeisTestConfig, notflag bool) {
	cmd := itutils.GetCommand("ps", "list")
	itutils.CheckList(t, params, cmd, "web.2 up (v2)", notflag)
}

func psScaleTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("ps", "scale")
	itutils.Execute(t, cmd, params, false, "")
}

func TestBuilds(t *testing.T) {
	params := psSetup(t)
	psScaleTest(t, params)
	appsOpenTest(t, params)
	psListTest(t, params, false)
	itutils.AppsDestroyTest(t, params)
	cmd := itutils.GetCommand("ps", "list")
	itutils.Execute(t, cmd, params, true, "404 NOT FOUND")
}
