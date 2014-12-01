// +build integration

package tests

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/deis/deis/tests/utils"
)

var (
	buildsListCmd   = "builds:list --app={{.AppName}}"
	buildsCreateCmd = `builds:create {{.ImageID}} --app={{.AppName}} --procfile="worker: while true; do echo hi; sleep 3; done"`
)

func TestBuilds(t *testing.T) {
	params := buildSetup(t)
	buildsListTest(t, params)
	appsOpenTest(t, params)
	utils.AppsDestroyTest(t, params)
	buildsCreateTest(t, params)
	// TODO: router needs a few seconds to wake up here--fixme!
	time.Sleep(5000 * time.Millisecond)
	appsOpenTest(t, params)
	buildsScaleTest(t, params)
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
	utils.Execute(t, "git commit --allow-empty -m bump", cfg, false, "")
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

// buildsCreateTest uses the `deis builds:create` (or `deis pull`) command
// to promote a build from an existing docker image.
func buildsCreateTest(t *testing.T, params *utils.DeisTestConfig) {
	params.AppName = "deispullsample"
	params.ImageID = "deis/example-go:latest"
	params.ExampleApp = "example-deis-pull"
	if err := os.Mkdir(params.ExampleApp, 0755); err != nil {
		t.Fatal(err)
	}
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, appsCreateCmdNoRemote, params, false, "")
	utils.Execute(t, buildsCreateCmd, params, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
}

// buildsScaleTest ensures that we can use a Procfile-based workflow for `deis pull`.
func buildsScaleTest(t *testing.T, params *utils.DeisTestConfig) {
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, "scale worker=1 --app={{.AppName}}", params, false, "")
	utils.Execute(t, "logs --app={{.AppName}}", params, false, "hi")
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
}
