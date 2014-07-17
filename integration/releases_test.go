package verbose

import (
	_ "fmt"
	"github.com/deis/deis/tests/integration-utils"
	"github.com/deis/deis/tests/utils"
	"testing"
)

func releasesSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	cfg.ExampleApp = itutils.GetRandomApp()
	cfg.AppName = "releasessample"
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
	cmd = itutils.GetCommand("config", "set")
	itutils.Execute(t, cmd, cfg, false, "")
	return cfg
}

func releasesListTest(t *testing.T, params *itutils.DeisTestConfig, notflag bool) {
	cmd := itutils.GetCommand("releases", "list")
	itutils.CheckList(t, params, cmd, params.Version, notflag)
}

func releasesInfoTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("releases", "info")
	itutils.Execute(t, cmd, params, false, "")
}

func releasesRollbackTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("releases", "rollback")
	itutils.Execute(t, cmd, params, false, "")
}

func appsOpenTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Curl(t, params)
}

func TestReleases(t *testing.T) {
	params := releasesSetup(t)
	releasesListTest(t, params, false)
	releasesInfoTest(t, params)
	releasesRollbackTest(t, params)
	appsOpenTest(t, params)
	params.Version = "4"
	releasesListTest(t, params, false)
	itutils.AppsDestroyTest(t, params)

}
