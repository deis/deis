// +build integration

package tests

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

// A Deis test configuration allows tests to be repeated against different
// targets, with different example apps, using specific credentials, and so on.
type deisTestConfig struct {
	AuthKey    string
	ExampleApp string
	Hosts      string
	HostName   string
	SshKey     string
}

// Test configuration created from environment variables (at compile time).
var envCfg = deisTestConfig{
	os.Getenv("AUTH_KEY"),
	os.Getenv("DEIS_TEST_APP"),
	os.Getenv("DEIS_TEST_HOSTS"),
	os.Getenv("DEIS_TEST_HOSTNAME"),
	os.Getenv("DEIS_TEST_SSH_KEY"),
}

// A test case is a relative directory plus a command that is expected to
// return 0 for success.
// The cmd field is run as an argument to "sh -c", so it can be arbitrarily
// complex.
type deisTest struct {
	dir string
	cmd string
}

// Tests to exercise a basic Deis workflow.
var vagrantTests = []deisTest{
	// Generate and activate a new SSH key named "deis".
	{"", `
if [ ! -f {{.AuthKey}} ]; then
  ssh-keygen -q -t rsa -f {{.AuthKey}} -N '' -C deis
  ssh-add {{.AuthKey}}
fi
`},
	// Register a "test" Deis user with the CLI, or skip if already registered.
	{"", `
deis register http://{{.HostName}}:8000 \
  --username=test \
  --password=asdf1234 \
  --email=test@test.co.nz || true
`},
	// Log in as the "test" user.
	{"", `
deis login http://{{.HostName}}:8000 \
  --username=test \
  --password=asdf1234
`},
	// Add the "deis" SSH key, or skip if it's been added already.
	{"", `
deis keys:add {{.AuthKey}}.pub || true
`},
	// Destroy the "dev" cluster if it exists.
	{"", `
deis clusters:destroy dev --confirm=dev || true
`},
	// Create a cluster named "dev".
	{"", `
deis init dev {{.HostName}} --hosts={{.Hosts}} --auth={{.SshKey}}
`},
	// Clone the example app git repository locally.
	{"", `
if [ ! -d ./{{.ExampleApp}} ]; then
  git clone https://github.com/deis/{{.ExampleApp}}.git
fi
`},
	// Remove the stale "deis" git remote if it exists.
	{"{{.ExampleApp}}", `
git remote remove deis || true
`},
	// TODO: GH issue about this sleep hack
	// Create an app named "testing".
	{"{{.ExampleApp}}", `
sleep 6 && deis apps:create testing
`},
	// git push the app to Deis
	{"{{.ExampleApp}}", `
git push deis master
`},
	// TODO: GH issue about this sleep hack
	// Test that the app's URL responds with "Powered by Deis".
	{"{{.ExampleApp}}", `
sleep 6 && curl -s http://testing.{{.HostName}} | grep -q 'Powered by Deis'
`},
	// Scale the app's web containers up to 3.
	{"{{.ExampleApp}}", `
deis scale web=3
`},
	// Test that the app's URL responds with "Powered by Deis".
	{"{{.ExampleApp}}", `
sleep 7 && curl -s http://testing.{{.HostName}} | grep -q 'Powered by Deis'
`},
}

// Updates a Vagrant instance to run Deis with docker containers using the
// current codebase, then registers a user, pushes an example app, and looks
// for "Powered by Deis" in the HTTP response.
func TestVagrantExampleApp(t *testing.T) {
	// t.Parallel()

	cfg := envCfg
	if cfg.AuthKey == "" {
		cfg.AuthKey = "~/.ssh/deis"
	}
	if cfg.ExampleApp == "" {
		cfg.ExampleApp = "example-ruby-sinatra"
	}
	if cfg.Hosts == "" {
		cfg.Hosts = "172.17.8.100"
	}
	if cfg.HostName == "" {
		cfg.HostName = "local.deisapp.com"
	}
	if cfg.SshKey == "" {
		cfg.SshKey = "~/.vagrant.d/insecure_private_key"
	}

	for _, tt := range vagrantTests {
		runTest(t, &tt, &cfg)
	}
}

var wd, _ = os.Getwd()

// Runs a test case and logs the results.
func runTest(t *testing.T, tt *deisTest, cfg *deisTestConfig) {
	// Fill in the command string template from our test configuration.
	var cmdBuf bytes.Buffer
	tmpl := template.Must(template.New("cmd").Parse(tt.cmd))
	if err := tmpl.Execute(&cmdBuf, cfg); err != nil {
		t.Fatal(err)
	}
	cmdString := cmdBuf.String()
	// Change to the target directory if needed.
	if tt.dir != "" {
		// Fill in the directory template from our test configuration.
		var dirBuf bytes.Buffer
		tmpl := template.Must(template.New("dir").Parse(tt.dir))
		if err := tmpl.Execute(&dirBuf, cfg); err != nil {
			t.Fatal(err)
		}
		dir, _ := filepath.Abs(filepath.Join(wd, dirBuf.String()))
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
	}
	// TODO: Go's testing package doesn't seem to allow for reporting interim
	// progress--we have to wait until everything completes (or fails) to see
	// anything that was written with t.Log or t.Fatal. Interim output would
	// be extremely helpful here, as this takes a while.
	// Execute the command and log the input and output on error.
	fmt.Printf("%v ... ", strings.TrimSpace(cmdString))
	cmd := exec.Command("sh", "-c", cmdString)
	if out, err := cmd.Output(); err != nil {
		t.Fatalf("%v\nOutput:\n%v", err, string(out))
	} else {
		fmt.Println("ok")
	}
}
