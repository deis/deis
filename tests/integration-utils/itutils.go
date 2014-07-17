package itutils

import (
	"bytes"
	"fmt"
	"github.com/ThomasRooney/gexpect"
	gson "github.com/bitly/go-simplejson"
	"github.com/deis/deis/tests/utils"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"text/template"
	"time"
)

var Deis = "/usr/local/bin/deis "

type DeisTestConfig struct {
	AuthKey      string
	Hosts        string
	HostName     string
	SshKey       string
	ClusterName  string
	UserName     string
	Password     string
	Email        string
	UpdatedHosts string
	ExampleApp   string
	AppName      string
	ProcessNum   string
	ImageId      string
	Version      string
	AppUser      string
}

func GetGlobalConfig() *DeisTestConfig {
	var envCfg = DeisTestConfig{
		"deis",
		"172.17.8.100",
		"local.deisapp.com",
		"~/.vagrant.d/insecure_private_key",
		"dev",
		"test",
		"asdf1234",
		"test@test.co.nz",
		"172.17.8.100",
		"example-go",
		"sample",
		"2",
		"buildtest",
		"2",
		"test1",
	}
	return &envCfg
}

//Tests example apps are running or not

func Curl(t *testing.T, params *DeisTestConfig) {
	url := "http://" + params.AppName + "." + params.HostName
	response, err := http.Get(url)
	if err != nil {
		t.Fatalf("not reachable:\n%v", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	if params.AppName == "example-python-django" {
		if !strings.Contains(string(body), "Powered by django") {
			t.Fatalf("App not started")
		}
	} else if !strings.Contains(string(body), "Powered by Deis") {
		t.Fatalf("App not started")
	}
}

//gexpect implementation of auth cancel

func AuthCancel(t *testing.T, params *DeisTestConfig) {
	fmt.Println("deis auth:cancel")
	child, err := gexpect.Spawn("/usr/local/bin/deis auth:cancel")
	if err != nil {
		t.Fatalf("command not started\n%v", err)
	}
	fmt.Println("username:")
	err = child.Expect("username:")
	if err != nil {
		t.Fatalf("expect username failed\n%v", err)
	}
	child.SendLine(params.UserName)
	fmt.Print("password:")
	err = child.Expect("password:")
	if err != nil {
		t.Fatalf("expect password failed\n%v", err)
	}
	child.SendLine(params.Password)
	err = child.ExpectRegex("(y/n)")
	if err != nil {
		t.Fatalf("expect cancel \n%v", err)
	}
	child.SendLine("y")
	err = child.Expect("Account cancelled")
	if err != nil {
		t.Fatalf("command executiuon failed\n%v", err)
	}
	child.Close()

}

/*CheckList takes config , command to execute and contain string and notflag .
*	Executes the command and checks if the contain string should be present or not according to notflag */

func CheckList(t *testing.T, params interface{}, cmd, contain string, notflag bool) {
	var cmdBuf bytes.Buffer
	tmpl := template.Must(template.New("cmd").Parse(cmd))
	if err := tmpl.Execute(&cmdBuf, params); err != nil {
		t.Fatal(err)
	}
	cmdString := cmdBuf.String()
	fmt.Println(cmdString)
	var cmdl *exec.Cmd
	if strings.Contains(cmd, "cat") {
		cmdl = exec.Command("sh", "-c", cmdString)
	} else {
		cmdl = exec.Command("sh", "-c", Deis+cmdString)
	}
	if stdout, _, err := utils.RunCommandWithStdoutStderr(cmdl); err == nil {
		if strings.Contains(stdout.String(), contain) != notflag {
			fmt.Println("Command Executed perfectly")
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
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
	var cmdl *exec.Cmd
	if strings.Contains(cmd, "git") {
		cmdl = exec.Command("sh", "-c", cmdString)
	} else {
		cmdl = exec.Command("sh", "-c", Deis+cmdString)
	}

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

//Destroys an app after execution of  each integration test

func AppsDestroyTest(t *testing.T, params *DeisTestConfig) {
	cmd := GetCommand("apps", "destroy")
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
	Execute(t, cmd, params, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
	if err := utils.Rmdir(params.ExampleApp); err != nil {
		t.Fatalf("Failed:\n%v", err)
	}
}

//Fetch commands from testconfig.json

func GetCommand(cmdtype, cmd string) string {
	js, _ := gson.NewJson(utils.GetFileBytes("testconfig.json"))
	command, _ := js.Get("commands").Get(cmdtype).Get(cmd).String()
	return command
}

//Selects a random app

func GetRandomApp() string {
	s1 := rand.NewSource(int64(time.Now().Unix()))
	r1 := rand.New(s1)
	appmap := make(map[int]string)
	appmap[0] = "example-go"
	appmap[1] = "example-ruby-sinatra"
	appmap[2] = "example-java-jetty"
	appmap[3] = "example-nodejs-express"
	appmap[4] = "example-python-flask"
	appmap[5] = "example-dockerfile-python"
	appmap[6] = "example-scala"
	appmap[7] = "example-clojure-ring"
	appmap[8] = "example-python-django"
	app := appmap[r1.Intn(8)]
	return app
}
