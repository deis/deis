package verbose

import (
	"bytes"
	"fmt"
	"github.com/deis/deis/tests/integration-utils"
	"github.com/deis/deis/tests/utils"
	"os/exec"
	"strings"
	"testing"
	"text/template"
)

func buildSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	cfg.ExampleApp = itutils.GetRandomApp()
	cfg.AppName = "buildsample"
	cmd := itutils.GetCommand("auth", "login")
	itutils.Execute(t, cmd, cfg, false, "")
	cmd = itutils.GetCommand("git", "clone")
	itutils.Execute(t, cmd, cfg, false, "")
	cmd = itutils.GetCommand("apps", "create")
	cmd1 := itutils.GetCommand("git", "push")
	cmd2 := itutils.GetCommand("git", "add")
	cmd3 := itutils.GetCommand("git", "commit")
	if err := utils.Chdir(cfg.ExampleApp); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}

	itutils.Execute(t, cmd, cfg, false, "")
	itutils.Execute(t, cmd1, cfg, false, "")
	if err := utils.CreateFile(cfg.ExampleApp); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
	itutils.Execute(t, cmd2, cfg, false, "")
	itutils.Execute(t, cmd3, cfg, false, "")
	itutils.Execute(t, cmd1, cfg, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
	return cfg
}

func buildsListTest(t *testing.T, params *itutils.DeisTestConfig) {
	Deis := "/usr/local/bin/deis "
	cmd := itutils.GetCommand("builds", "list")
	var cmdBuf bytes.Buffer
	tmpl := template.Must(template.New("cmd").Parse(cmd))
	if err := tmpl.Execute(&cmdBuf, params); err != nil {
		t.Fatal(err)
	}
	cmdString := cmdBuf.String()
	fmt.Println(cmdString)
	cmdl := exec.Command("sh", "-c", Deis+cmdString)
	if stdout, _, err := utils.RunCommandWithStdoutStderr(cmdl); err != nil {
		t.Fatalf("Failed:\n%v", err)
	} else {
		ImageId := strings.Split(stdout.String(), "\n")[2]
		params.ImageId = strings.Fields(ImageId)[0]
	}

}

func buildsCreateTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("builds", "create")
	itutils.Execute(t, cmd, params, false, "")

}

func TestBuilds(t *testing.T) {
	params := buildSetup(t)
	buildsListTest(t, params)
	buildsCreateTest(t, params)
	itutils.AppsDestroyTest(t, params)
}
