// +build integration

package tests

import (
	"testing"
	"time"

	"github.com/deis/deis/tests/integration-utils"
	"github.com/deis/deis/tests/utils"
)

var (
	appsCreateCmd  = "apps:create {{.AppName}}"
	appsListCmd    = "apps:list"
	appsRunCmd     = "apps:run echo hello"
	appsOpenCmd    = "apps:open --app={{.AppName}}"
	appsLogsCmd    = "apps:logs --app={{.AppName}}"
	appsInfoCmd    = "apps:info --app={{.AppName}}"
	appsDestroyCmd = "apps:destroy --app={{.AppName}} --confirm={{.AppName}}"
)

func appsSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	cfg.AppName = "appssample"
	itutils.Execute(t, authLoginCmd, cfg, false, "")
	itutils.Execute(t, gitCloneCmd, cfg, false, "")
	return cfg
}

func appsCreateTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := appsCreateCmd
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	itutils.Execute(t, cmd, params, false, "")
	itutils.Execute(t, cmd, params, true, "App with this Id already exists")
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
}

func appsRunTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := appsRunCmd
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	itutils.Execute(t, cmd, params, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	itutils.Execute(t, cmd, params, true, "Not found")
}

func appsDestroyTest(t *testing.T, params *itutils.DeisTestConfig) {
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	itutils.Execute(t, appsDestroyCmd, params, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	if err := utils.Rmdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
}

func appsListTest(t *testing.T, params *itutils.DeisTestConfig, notflag bool) {
	itutils.CheckList(t, params, appsListCmd, params.AppName, notflag)
}

func appsLogsTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := appsLogsCmd
	itutils.Execute(t, cmd, params, true, "204 NO CONTENT")
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	itutils.Execute(t, gitPushCmd, params, false, "")
	// TODO: nginx needs a few seconds to wake up here--fixme!
	time.Sleep(5000 * time.Millisecond)
	itutils.Curl(t, params)
	itutils.Execute(t, cmd, params, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
}

func appsInfoTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Execute(t, appsInfoCmd, params, false, "")
}

func appsOpenTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Curl(t, params)
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
