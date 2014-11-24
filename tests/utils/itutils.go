package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/ThomasRooney/gexpect"
)

// Deis points to the CLI used to run tests.
var Deis = "deis "

// DeisTestConfig allows tests to be repeated against different
// targets, with different example apps, using specific credentials, and so on.
type DeisTestConfig struct {
	AuthKey     string
	Hosts       string
	Domain      string
	SSHKey      string
	ClusterName string
	UserName    string
	Password    string
	Email       string
	ExampleApp  string
	AppName     string
	ProcessNum  string
	ImageID     string
	Version     string
	AppUser     string
}

// randomApp is used for the test run if DEIS_TEST_APP isn't set
var randomApp = GetRandomApp()

// GetGlobalConfig returns a test configuration object.
func GetGlobalConfig() *DeisTestConfig {
	authKey := os.Getenv("DEIS_TEST_AUTH_KEY")
	if authKey == "" {
		authKey = "deis"
	}
	hosts := os.Getenv("DEIS_TEST_HOSTS")
	if hosts == "" {
		hosts = "172.17.8.100"
	}
	domain := os.Getenv("DEIS_TEST_DOMAIN")
	if domain == "" {
		domain = "local.deisapp.com"
	}
	sshKey := os.Getenv("DEIS_TEST_SSH_KEY")
	if sshKey == "" {
		sshKey = "~/.vagrant.d/insecure_private_key"
	}
	exampleApp := os.Getenv("DEIS_TEST_APP")
	if exampleApp == "" {
		exampleApp = randomApp
	}
	var envCfg = DeisTestConfig{
		AuthKey:     authKey,
		Hosts:       hosts,
		Domain:      domain,
		SSHKey:      sshKey,
		ClusterName: "dev",
		UserName:    "test",
		Password:    "asdf1234",
		Email:       "test@test.co.nz",
		ExampleApp:  exampleApp,
		AppName:     "sample",
		ProcessNum:  "2",
		ImageID:     "buildtest",
		Version:     "2",
		AppUser:     "test1",
	}
	return &envCfg
}

func doCurl(url string) ([]byte, error) {
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)

	if !strings.Contains(string(body), "Powered by Deis") {
		return nil, fmt.Errorf("App not started (%d)\nBody: (%s)", response.StatusCode, string(body))
	}

	return body, nil
}

// Curl connects to a Deis endpoint to see if the example app is running.
func Curl(t *testing.T, params *DeisTestConfig) {
	CurlWithFail(t, params, false, "")
}

// Curl connects to a Deis endpoint to see if the example app is running.
func CurlWithFail(t *testing.T, params *DeisTestConfig, failFlag bool, expect string) {
	url := "http://" + params.AppName + "." + params.Domain

	// FIXME: try the curl a few times
	for i := 0; i < 20; i++ {
		body, err := doCurl(url)
		if err == nil {
			fmt.Println(string(body))
			return
		}
		time.Sleep(1 * time.Second)
	}

	// once more to fail with an error
	body, err := doCurl(url)

	switch failFlag {
	case true:
		if err != nil {
			if strings.Contains(string(err.Error()), expect) {
				fmt.Println("(Error expected...ok) " + expect)
			} else {
				t.Fatal(err)
			}
		} else {
			if strings.Contains(string(body), expect) {
				fmt.Println("(Error expected...ok) " + expect)
			} else {
				t.Fatal(err)
			}
		}
	case false:
		if err != nil {
			t.Fatal(err)
		} else {
			fmt.Println(string(body))
		}
	}
}

// AuthCancel tests whether `deis auth:cancel` destroys a user's account.
func AuthCancel(t *testing.T, params *DeisTestConfig) {
	fmt.Println("deis auth:cancel")
	child, err := gexpect.Spawn(Deis + " auth:cancel")
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

// AuthPasswd tests whether `deis auth:passwd` updates a user's password.
func AuthPasswd(t *testing.T, params *DeisTestConfig, password string) {
	fmt.Println("deis auth:passwd")
	child, err := gexpect.Spawn(Deis + " auth:passwd")
	if err != nil {
		t.Fatalf("command not started\n%v", err)
	}
	fmt.Println("current password:")
	err = child.Expect("current password: ")
	if err != nil {
		t.Fatalf("expect password failed\n%v", err)
	}
	child.SendLine(params.Password)
	fmt.Println("new password:")
	err = child.Expect("new password: ")
	if err != nil {
		t.Fatalf("expect password failed\n%v", err)
	}
	child.SendLine(password)
	fmt.Println("new password (confirm):")
	err = child.Expect("new password (confirm): ")
	if err != nil {
		t.Fatalf("expect password failed\n%v", err)
	}
	child.SendLine(password)
	err = child.Expect("Password change succeeded")
	if err != nil {
		t.Fatalf("command executiuon failed\n%v", err)
	}
	child.Close()
}

// CheckList executes a command and optionally tests whether its output does
// or does not contain a given string.
func CheckList(
	t *testing.T, cmd string, params interface{}, contain string, notflag bool) {
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
	stdout, _, err := RunCommandWithStdoutStderr(cmdl)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(stdout.String(), contain) == notflag {
		if notflag {
			t.Fatalf(
				"Didn't expect '%s' in command output:\n%v", contain, stdout)
		}
		t.Fatalf("Expected '%s' in command output:\n%v", contain, stdout)
	}
}

// Execute takes command string and parameters required to execute the command,
// a failflag to check whether the command is expected to fail, and an expect
// string to check whether the command has failed according to failflag.
//
// If failflag is true and the command failed, check the stdout and stderr for
// the expect string.
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
		if stdout, stderr, err := RunCommandWithStdoutStderr(cmdl); err != nil {
			if strings.Contains(stdout.String(), expect) || strings.Contains(stderr.String(), expect) {
				fmt.Println("(Error expected...ok)")
			} else {
				t.Fatal(err)
			}
		} else {
			if strings.Contains(stdout.String(), expect) || strings.Contains(stderr.String(), expect) {
				fmt.Println("(Error expected...ok)" + expect)
			} else {
				t.Fatal(err)
			}
		}
	case false:
		if _, _, err := RunCommandWithStdoutStderr(cmdl); err != nil {
			t.Fatal(err)
		} else {
			fmt.Println("ok")
		}
	}
}

// AppsDestroyTest destroys a Deis app and checks that it was successful.
func AppsDestroyTest(t *testing.T, params *DeisTestConfig) {
	cmd := "apps:destroy --app={{.AppName}} --confirm={{.AppName}}"
	if err := Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	Execute(t, cmd, params, false, "")
	if err := Chdir(".."); err != nil {
		t.Fatal(err)
	}
	if err := Rmdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
}

// GetRandomApp returns a known working example app at random for testing.
func GetRandomApp() string {
	rand.Seed(int64(time.Now().Unix()))
	apps := []string{
		"example-clojure-ring",
		// "example-dart",
		"example-dockerfile-python",
		"example-go",
		"example-java-jetty",
		"example-nodejs-express",
		// "example-php",
		"example-play",
		"example-python-django",
		"example-python-flask",
		"example-ruby-sinatra",
		"example-scala",
	}
	return apps[rand.Intn(len(apps))]
}
