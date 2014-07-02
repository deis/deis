package itutils

import (
	"bytes"
	"fmt"
	gson "github.com/bitly/go-simplejson"
	"github.com/deis/deis/tests/utils"
	"os/exec"
	"strings"
	"testing"
	"text/template"
)

var Deis = "/usr/local/bin/deis "

type DeisTestConfig struct {
	AuthKey  string
	Hosts    string
	HostName string
	SshKey   string
}

type UserDetails struct {
	UserName string
	Password string
	Email    string
	HostName string
}

func SetUser() *UserDetails {
	var user = new(UserDetails)
	user.UserName, user.Password = utils.GetUserDetails()
	user.Email = "test@test.co.nz"
	user.HostName = "deis.54.193.35.8.xip.io"
	return user
}

func GlobalSetup(t *testing.T) *DeisTestConfig {
	var envCfg = DeisTestConfig{
		"~/.ssh/deis",
		"54.193.35.8",
		"deis.54.193.35.8.xip.io",
		"~/.vagrant.d/insecure_private_key",
	}
	var user = UserDetails{
		"Test",
		"asdf1234",
		"test@test.co.nz",
		envCfg.HostName,
	}
	Execute(t, GetCommand("auth", "register"), user, false, "")
	return &envCfg
}

/***Execute function takes command string and parameters required to execute the command
A failflag to check whether the command is expected to fail
An expect string to check whether the command has failed according to failflag

If a failflag is true and command failed we check the stdout and stderr for expect string

***/


func Execute(t *testing.T, cmd string, params interface{}, failFlag bool, expect string) {
	var cmdBuf bytes.Buffer
	tmpl := template.Must(template.New("cmd").Parse(cmd))
	if err := tmpl.Execute(&cmdBuf, params); err != nil {
		t.Fatal(err)
	}
	cmdString := cmdBuf.String()
	fmt.Println(cmdString)
	cmdl := exec.Command("sh", "-c", Deis+cmdString)

	switch failFlag {
	case true:
		if stdout, stderr, err := utils.RunCommandWithStdoutStderr(cmdl); err != nil {
			if strings.Contains(stdout.String(), expect) || strings.Contains(stderr.String(), expect) {
				fmt.Println("Test Failed Expected behavior")
			} else {
				t.Fatalf("Failed:\n%v", err)
			}
		} else {
			if strings.Contains(stdout.String(), expect) || strings.Contains(stderr.String(), expect) {
				fmt.Println("expected" + expect)
			} else {
				t.Fatalf("Failed:\n%v", err)
			}
		}
	case false:
		if _, _, err := utils.RunCommandWithStdoutStderr(cmdl); err != nil {
			t.Fatalf("Failed:\n%v", err)
		} else {
			fmt.Println("ok")
		}
	}
}



func GetCommand(cmdtype, cmd string) string {
	js, _ := gson.NewJson(utils.GetFileBytes("testconfig.json"))
	command, _ := js.Get("commands").Get(cmdtype).Get(cmd).String()
	return command
}
