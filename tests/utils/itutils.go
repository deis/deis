package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
	AuthKey            string
	Hosts              string
	Domain             string
	SSHKey             string
	ClusterName        string
	UserName           string
	Password           string
	Email              string
	ExampleApp         string
	AppDomain          string
	AppName            string
	ProcessNum         string
	ImageID            string
	Version            string
	AppUser            string
	SSLCertificatePath string
	SSLKeyPath         string
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
	appDomain := os.Getenv("DEIS_TEST_APP_DOMAIN")
	if appDomain == "" {
		appDomain = fmt.Sprintf("test.%s", domain)
	}

	// generate a self-signed certifcate for the app domain
	keyOut, err := filepath.Abs(appDomain + ".key")
	if err != nil {
		log.Fatal(err)
	}
	certOut, err := filepath.Abs(appDomain + ".cert")
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("openssl", "req", "-new", "-newkey", "rsa:4096", "-nodes", "-x509",
		"-days", "1",
		"-subj", fmt.Sprintf("/C=US/ST=Colorado/L=Boulder/CN=%s", appDomain),
		"-keyout", keyOut,
		"-out", certOut)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	var envCfg = DeisTestConfig{
		AuthKey:            authKey,
		Hosts:              hosts,
		Domain:             domain,
		SSHKey:             sshKey,
		ClusterName:        "dev",
		UserName:           "test",
		Password:           "asdf1234",
		Email:              "test@test.co.nz",
		ExampleApp:         exampleApp,
		AppDomain:          appDomain,
		AppName:            "sample",
		ProcessNum:         "2",
		ImageID:            "buildtest",
		Version:            "2",
		AppUser:            "test1",
		SSLCertificatePath: certOut,
		SSLKeyPath:         keyOut,
	}
	return &envCfg
}

func doCurl(url string) ([]byte, error) {
	// disable security check for self-signed certificates
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if !strings.Contains(string(body), "Powered by") {
		return nil, fmt.Errorf("App not started (%d)\nBody: (%s)", response.StatusCode, string(body))
	}

	return body, nil
}

// Curl connects to an endpoint to see if the endpoint is responding.
func Curl(t *testing.T, url string) {
	CurlWithFail(t, url, false, "")
}

// CurlApp is a convenience function to see if the example app is running.
func CurlApp(t *testing.T, cfg DeisTestConfig) {
	CurlWithFail(t, fmt.Sprintf("http://%s.%s", cfg.AppName, cfg.Domain), false, "")
}

// CurlWithFail connects to a Deis endpoint to see if the example app is running.
func CurlWithFail(t *testing.T, url string, failFlag bool, expect string) {
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
				if expect == "503" {
					if strings.Contains(string(err.Error()), "Welcome to nginx!") {
						return
					}
				}
				t.Fatal(err)
			}
		} else {
			if strings.Contains(string(body), expect) {
				fmt.Println("(Error expected...ok) " + expect)
			} else {
				if expect == "503" {
					if strings.Contains(string(err.Error()), "Welcome to nginx!") {
						return
					}
				}
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
	if notflag && strings.Contains(stdout.String(), contain) {
		t.Fatalf("Didn't expect '%s' in command output:\n%v", contain, stdout)
	}
	if !notflag && !strings.Contains(stdout.String(), contain) {
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
	if strings.Contains(cmd, "git ") {
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
		"example-dockerfile-http",
	}
	return apps[rand.Intn(len(apps))]
}
