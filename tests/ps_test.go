// +build integration

package tests

import (
	_ "fmt"
	"strings"
	"testing"

	"github.com/deis/deis/tests/integration-utils"
	"github.com/deis/deis/tests/utils"
)

func psSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
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
	output := "web.2 up (v2)"
	if strings.Contains(params.ExampleApp, "dockerfile") {
		output = strings.Replace(output, "web", "cmd", 1)
	}
	itutils.CheckList(t, params, cmd, output, notflag)
}

func psScaleTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("ps", "scale")
	if strings.Contains(params.ExampleApp, "dockerfile") {
		cmd = strings.Replace(cmd, "web", "cmd", 1)
	}
	itutils.Execute(t, cmd, params, false, "")
}

func TestPs(t *testing.T) {
	params := psSetup(t)
	psScaleTest(t, params)
	appsOpenTest(t, params)
	psListTest(t, params, false)
	itutils.AppsDestroyTest(t, params)
	cmd := itutils.GetCommand("ps", "list")
	itutils.Execute(t, cmd, params, true, "404 NOT FOUND")
}
