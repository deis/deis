// +build integration

package tests

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
	"text/template"

	"github.com/deis/deis/tests/utils"
)

var (
	buildsListCmd   = "builds:list --app={{.AppName}}"
	buildsCreateCmd = "builds:create {{.ImageID}} --app={{.AppName}}"
)

func TestBuilds(t *testing.T) {
	params := buildSetup(t)
	buildsListTest(t, params)
	buildsCreateTest(t, params)
	appsOpenTest(t, params)
	utils.AppsDestroyTest(t, params)
}

func buildSetup(t *testing.T) *utils.DeisTestConfig {
	cfg := utils.GetGlobalConfig()
	cfg.AppName = "buildsample"
	utils.Execute(t, authLoginCmd, cfg, false, "")
	utils.Execute(t, gitCloneCmd, cfg, false, "")
	if err := utils.Chdir(cfg.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, appsCreateCmd, cfg, false, "")
	utils.Execute(t, gitPushCmd, cfg, false, "")
	if err := utils.CreateFile(cfg.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, gitAddCmd, cfg, false, "")
	utils.Execute(t, gitCommitCmd, cfg, false, "")
	utils.Execute(t, gitPushCmd, cfg, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	return cfg
}

func buildsListTest(t *testing.T, params *utils.DeisTestConfig) {
	cmd := buildsListCmd
	var cmdBuf bytes.Buffer
	tmpl := template.Must(template.New("cmd").Parse(cmd))
	if err := tmpl.Execute(&cmdBuf, params); err != nil {
		t.Fatal(err)
	}
	cmdString := cmdBuf.String()
	cmdl := exec.Command("sh", "-c", utils.Deis+cmdString)
	stdout, _, err := utils.RunCommandWithStdoutStderr(cmdl)
	if err != nil {
		t.Fatal(err)
	}
	ImageID := strings.Split(stdout.String(), "\n")[2]
	params.ImageID = strings.Fields(ImageID)[0]
}

func buildsCreateTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.Execute(t, buildsCreateCmd, params, false, "")
}
