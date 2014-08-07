// +build integration

package tests

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/deis/deis/tests/integration-utils"
	"github.com/deis/deis/tests/utils"
)

var (
	psListCmd  = "ps:list --app={{.AppName}}"
	psScaleCmd = "ps:scale web={{.ProcessNum}} --app={{.AppName}}"
)

func psSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	cfg.AppName = "pssample"
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

func psListTest(t *testing.T, params *itutils.DeisTestConfig, notflag bool) {
	output := "web.2 up (v2)"
	if strings.Contains(params.ExampleApp, "dockerfile") {
		output = strings.Replace(output, "web", "cmd", 1)
	}
	itutils.CheckList(t, params, psListCmd, output, notflag)
}

func psScaleTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := psScaleCmd
	if strings.Contains(params.ExampleApp, "dockerfile") {
		cmd = strings.Replace(cmd, "web", "cmd", 1)
	}
	itutils.Execute(t, cmd, params, false, "")
	// Regression test for https://github.com/deis/deis/pull/1347
	// Ensure that systemd unitfile droppings are cleaned up.
	sshCmd := exec.Command("ssh",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "PasswordAuthentication=no",
		"core@deis."+params.Domain, "ls")
	out, err := sshCmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), ".service") {
		t.Fatalf("systemd files left on filesystem: \n%s", out)
	}
}

func TestPs(t *testing.T) {
	params := psSetup(t)
	psScaleTest(t, params)
	appsOpenTest(t, params)
	psListTest(t, params, false)
	itutils.AppsDestroyTest(t, params)
	itutils.Execute(t, psScaleCmd, params, true, "404 NOT FOUND")
}
