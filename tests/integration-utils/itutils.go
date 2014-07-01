package itutils

import (
	"bytes"
	"fmt"
	gson "github.com/bitly/go-simplejson"
	"github.com/deis/deis/tests/utils"
	"os/exec"
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
	return user
}

func GlobalSetup(t *testing.T) *DeisTestConfig {
	var envCfg = DeisTestConfig{
		"~/.ssh/deis",
		"54.193.41.120",
		"deis.54.193.9.175.xip.io",
		"~/.vagrant.d/insecure_private_key",
	}
	var user = UserDetails{
		"Test",
		"asdf1234",
		"test@test.co.nz",
		envCfg.HostName,
	}
	Execute(t, GetCommand("auth", "register"), user, false)
	return &envCfg
}

func Execute(t *testing.T, cmd string, params interface{}, failFlag bool, expect string) {
	var cmdBuf bytes.Buffer
	tmpl := template.Must(template.New("cmd").Parse(cmd))
	if err := tmpl.Execute(&cmdBuf, params); err != nil {
		t.Fatal(err)
	}
	cmdString := cmdBuf.String()
	fmt.Println(cmdString)
	cmdl := exec.Command("sh", "-c", Deis+cmdString)
	if err := utils.RunCommandWithStdoutStderr(cmdl); err != nil {
		if failFlag {
			fmt.Println("Test Failed expected behavior ")
		} else {
			t.Fatalf("Output:\n%v", err)
		}
	} else if failFlag {
		t.Fatalf("test should be failed here but passing ")
	} else {
		fmt.Println("ok")
	}
}

func GetCommand(cmdtype, cmd string) string {
	js, _ := gson.NewJson(utils.GetFileBytes("testconfig.json"))
	command, _ := js.Get("commands").Get(cmdtype).Get(cmd).String()
	return command
}
