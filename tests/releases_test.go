// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/utils"
)

var (
	releasesListCmd     = "releases:list --app={{.AppName}}"
	releasesInfoCmd     = "releases:info {{.Version}} --app={{.AppName}}"
	releasesRollbackCmd = "releases:rollback {{.Version}} --app={{.AppName}}"
)

func TestReleases(t *testing.T) {
	params := releasesSetup(t)
	releasesListTest(t, params, false)
	releasesInfoTest(t, params)
	releasesRollbackTest(t, params)
	appsOpenTest(t, params)
	params.Version = "4"
	releasesListTest(t, params, false)
	utils.AppsDestroyTest(t, params)

}

func releasesSetup(t *testing.T) *utils.DeisTestConfig {
	cfg := utils.GetGlobalConfig()
	cfg.AppName = "releasessample"
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
	utils.Execute(t, configSetCmd, cfg, false, "")
	return cfg
}

func releasesInfoTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.Execute(t, releasesInfoCmd, params, false, "")
}

func releasesListTest(
	t *testing.T, params *utils.DeisTestConfig, notflag bool) {
	utils.CheckList(t, releasesListCmd, params, params.Version, notflag)
}

func releasesRollbackTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.Execute(t, releasesRollbackCmd, params, false, "")
}
