// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/utils"
)

var (
	configListCmd  = "config:list --app={{.AppName}}"
	configSetCmd   = "config:set FOO=讲台 --app={{.AppName}}"
	configSet2Cmd  = "config:set FOO=10 --app={{.AppName}}"
	configUnsetCmd = "config:unset FOO --app={{.AppName}}"
)

func TestConfig(t *testing.T) {
	params := configSetup(t)
	configSetTest(t, params)
	configListTest(t, params, false)
	appsOpenTest(t, params)
	configUnsetTest(t, params)
	configListTest(t, params, true)
	limitsSetTest(t, params, 4)
	appsOpenTest(t, params)
	limitsUnsetTest(t, params, 6)
	appsOpenTest(t, params)
	//tagsTest(t, params, 8)
	appsOpenTest(t, params)
	utils.AppsDestroyTest(t, params)
}

func configSetup(t *testing.T) *utils.DeisTestConfig {
	cfg := utils.GetGlobalConfig()
	cfg.AppName = "configsample"
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

func configListTest(
	t *testing.T, params *utils.DeisTestConfig, notflag bool) {
	utils.CheckList(t, configListCmd, params, "FOO", notflag)
}

func configSetTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.Execute(t, configSetCmd, params, false, "讲台")
	utils.CheckList(t, appsInfoCmd, params, "(v3)", false)
	utils.Execute(t, configSet2Cmd, params, false, "10")
	utils.CheckList(t, appsInfoCmd, params, "(v4)", false)
}

func configUnsetTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.Execute(t, configUnsetCmd, params, false, "")
	utils.CheckList(t, appsInfoCmd, params, "(v5)", false)
	utils.CheckList(t, "run env --app={{.AppName}}", params, "FOO", true)
}
