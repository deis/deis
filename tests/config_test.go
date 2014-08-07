// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/integration-utils"
	"github.com/deis/deis/tests/utils"
)

var (
	configListCmd  = "config:list --app={{.AppName}}"
	configSetCmd   = "config:set jaf=1 --app={{.AppName}}"
	configUnsetCmd = "config:unset jaf --app={{.AppName}}"
)

func TestConfig(t *testing.T) {
	params := configSetup(t)
	configSetTest(t, params)
	configListTest(t, params, false)
	appsOpenTest(t, params)
	configUnsetTest(t, params)
	configListTest(t, params, true)
	itutils.AppsDestroyTest(t, params)
}

func configSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	cfg.AppName = "configsample"
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

func configListTest(
	t *testing.T, params *itutils.DeisTestConfig, notflag bool) {
	itutils.CheckList(t, configListCmd, params, "jaf", notflag)
}

func configSetTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Execute(t, configSetCmd, params, false, "")
	itutils.CheckList(t, appsInfoCmd, params, "(v3)", false)
}

func configUnsetTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Execute(t, configUnsetCmd, params, false, "")
	itutils.CheckList(t, appsInfoCmd, params, "(v4)", false)
}
