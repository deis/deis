// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/integration-utils"
	"github.com/deis/deis/tests/utils"
)

var (
	releasesListCmd     = "releases:list --app={{.AppName}}"
	releasesInfoCmd     = "releases:info {{.Version}} --app={{.AppName}}"
	releasesRollbackCmd = "releases:rollback {{.Version}} --app={{.AppName}}"
)

func releasesSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	cfg.AppName = "releasessample"
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
	itutils.Execute(t, configSetCmd, cfg, false, "")
	return cfg
}

func releasesListTest(
	t *testing.T, params *itutils.DeisTestConfig, notflag bool) {
	itutils.CheckList(t, params, releasesListCmd, params.Version, notflag)
}

func releasesInfoTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Execute(t, releasesInfoCmd, params, false, "")
}

func releasesRollbackTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Execute(t, releasesRollbackCmd, params, false, "")
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
