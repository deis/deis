// +build integration

package tests

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/deis/deis/tests/utils"
)

var (
	appsCreateCmd          = "apps:create {{.AppName}}"
	appsCreateCmdNoRemote  = "apps:create {{.AppName}} --no-remote"
	appsCreateCmdBuildpack = "apps:create {{.AppName}} --buildpack https://example.com"
	appsListCmd            = "apps:list"
	appsRunCmd             = "apps:run echo Hello, 世界"
	appsOpenCmd            = "apps:open --app={{.AppName}}"
	appsLogsCmd            = "apps:logs --app={{.AppName}}"
	appsLogsLimitCmd       = "apps:logs --app={{.AppName}} -n 1"
	appsInfoCmd            = "apps:info --app={{.AppName}}"
	appsDestroyCmd         = "apps:destroy --app={{.AppName}} --confirm={{.AppName}}"
	appsDestroyCmdNoApp    = "apps:destroy --confirm={{.AppName}}"
)

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
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

func appsSetup(t *testing.T) *utils.DeisTestConfig {
	cfg := utils.GetGlobalConfig()
	cfg.AppName = "appssample"
	utils.Execute(t, authLoginCmd, cfg, false, "")
	utils.Execute(t, gitCloneCmd, cfg, false, "")
	return cfg
}

func appsCreateTest(t *testing.T, params *utils.DeisTestConfig) {
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	// TODO: move --buildpack to client unit tests
	utils.Execute(t, appsCreateCmdBuildpack, params, false, "BUILDPACK_URL")
	utils.Execute(t, appsDestroyCmdNoApp, params, false, "")
	utils.Execute(t, appsCreateCmd, params, false, "")
	utils.Execute(t, appsCreateCmd, params, true, "This field must be unique.")
}

func appsDestroyTest(t *testing.T, params *utils.DeisTestConfig) {
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, appsDestroyCmd, params, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	if err := utils.Rmdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
}

func appsInfoTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.Execute(t, appsInfoCmd, params, false, "")
}

func appsListTest(t *testing.T, params *utils.DeisTestConfig, notflag bool) {
	utils.CheckList(t, appsListCmd, params, params.AppName, notflag)
}

func appsLogsTest(t *testing.T, params *utils.DeisTestConfig) {
	cmd := appsLogsCmd
	// test for application lifecycle logs
	utils.Execute(t, cmd, params, false, "204 NO CONTENT")
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, gitPushCmd, params, false, "")
	utils.CurlApp(t, *params)
	utils.Execute(t, cmd, params, false, "created initial release")
	utils.Execute(t, cmd, params, false, "listening on 5000...")

	utils.Execute(t, appsLogsLimitCmd, params, false, "")

	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
}

func appsOpenTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.CurlApp(t, *params)
	utils.CurlWithFail(t, fmt.Sprintf("http://%s.%s", "this-app-does-not-exist", params.Domain), true, "404 Not Found")
}

func appsRunTest(t *testing.T, params *utils.DeisTestConfig) {
	cmd := appsRunCmd
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.CheckList(t, cmd, params, "Hello, 世界", false)
	utils.Execute(t, "apps:run env", params, true, "GIT_SHA")
	// run a REALLY large command to test https://github.com/deis/deis/issues/2046
	largeString := randomString(1024)
	utils.Execute(t, "apps:run echo "+largeString, params, false, largeString)
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, cmd, params, true, "Not found")
}
