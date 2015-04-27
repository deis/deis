// +build integration

package tests

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/deis/deis/tests/utils"
)

var (
	psListCmd      = "ps:list --app={{.AppName}}"
	psScaleCmd     = "ps:scale web={{.ProcessNum}} --app={{.AppName}}"
	psDownScaleCmd = "ps:scale web=0 --app={{.AppName}}"
	psRestartCmd   = "ps:restart web --app={{.AppName}}"
)

func TestPs(t *testing.T) {
	params := psSetup(t)
	psScaleTest(t, params, psScaleCmd)
	appsOpenTest(t, params)
	psListTest(t, params, false)
	psScaleTest(t, params, psRestartCmd)
	psScaleTest(t, params, psDownScaleCmd)

	// FIXME if we don't wait here, some of the routers may give us a 502 before
	// the app is removed from the config.
	// we wait 7 seconds since confd reloads every 5 seconds
	time.Sleep(time.Millisecond * 7000)

	// test for a 503 response
	utils.CurlWithFail(t, fmt.Sprintf("http://%s.%s", params.AppName, params.Domain), true, "503")

	utils.AppsDestroyTest(t, params)
	utils.Execute(t, psScaleCmd, params, true, "404 NOT FOUND")
	// ensure we can choose our preferred beverage
	utils.Execute(t, psScaleCmd, params, true, "but first, coffee!")
	if err := os.Setenv("DEIS_DRINK_OF_CHOICE", "tea"); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, psScaleCmd, params, true, "but first, tea!")
}

func psSetup(t *testing.T) *utils.DeisTestConfig {
	cfg := utils.GetGlobalConfig()
	cfg.AppName = "pssample"
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

func psListTest(t *testing.T, params *utils.DeisTestConfig, notflag bool) {
	output := "web.2 up (v2)"
	if strings.Contains(params.ExampleApp, "dockerfile") {
		output = strings.Replace(output, "web", "cmd", 1)
	}
	utils.CheckList(t, psListCmd, params, output, notflag)
}

func psScaleTest(t *testing.T, params *utils.DeisTestConfig, cmd string) {
	if strings.Contains(params.ExampleApp, "dockerfile") {
		cmd = strings.Replace(cmd, "web", "cmd", 1)
	}
	utils.Execute(t, cmd, params, false, "")
}
