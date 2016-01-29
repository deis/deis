// +build integration

package tests

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/deis/deis/tests/utils"
)

var (
	appsCreateCmd          = "apps:create {{.AppName}}"
	appsCreateCmdNoRemote  = "apps:create {{.AppName}} --no-remote"
	appsCreateCmdBuildpack = "apps:create {{.AppName}} --buildpack https://example.com"
	appsListCmd            = "apps:list"
	appsRunCmd             = "apps:run echo Hello, 世界"
	appsOpenCmd            = "apps:open --app={{.AppName}}"
	appsLogsCmd            = "apps:logs --app={{.AppName}}"
	appsLogsLimitCmd       = "apps:logs --app={{.AppName}} -n 1"
	appsInfoCmd            = "apps:info --app={{.AppName}}"
	appsDestroyCmd         = "apps:destroy --app={{.AppName}} --confirm={{.AppName}}"
	appsDestroyCmdNoApp    = "apps:destroy --confirm={{.AppName}}"
	appsTransferCmd        = "apps:transfer {{.NewOwner}} --app={{.AppName}}"
)

func randomString(n int) string {
	// Be sure we've seeded the random number generator, otherwise we could get the same string
	// every time.
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestApps(t *testing.T) {
	params := appsSetup(t)
	appsCreateTest(t, params)
	appsListTest(t, params, false)
	appsLogsTest(t, params)
	appsInfoTest(t, params)
	appsRunTest(t, params)
	appsOpenTest(t, params)
	appsDestroyTest(t, params)
	appsListTest(t, params, true)
	appsTransferTest(t, params)
}

func appsSetup(t *testing.T) *utils.DeisTestConfig {
	cfg := utils.GetGlobalConfig()
	cfg.AppName = "appssample"
	utils.Execute(t, authLoginCmd, cfg, false, "")
	utils.Execute(t, gitCloneCmd, cfg, false, "")
	return cfg
}

func appsCreateTest(t *testing.T, params *utils.DeisTestConfig) {
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	// TODO: move --buildpack to client unit tests
	utils.Execute(t, appsCreateCmdBuildpack, params, false, "BUILDPACK_URL")
	utils.Execute(t, appsDestroyCmdNoApp, params, false, "")
	utils.Execute(t, appsCreateCmd, params, false, "")
	utils.Execute(t, appsCreateCmd, params, true, "This field must be unique.")
}

func appsDestroyTest(t *testing.T, params *utils.DeisTestConfig) {
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, appsDestroyCmd, params, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	if err := utils.Rmdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
}

func appsInfoTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.Execute(t, appsInfoCmd, params, false, "")
}

func appsListTest(t *testing.T, params *utils.DeisTestConfig, notflag bool) {
	utils.CheckList(t, appsListCmd, params, params.AppName, notflag)
}

func appsLogsTest(t *testing.T, params *utils.DeisTestConfig) {
	cmd := appsLogsCmd
	// test for application lifecycle logs
	utils.Execute(t, cmd, params, false, "204 NO CONTENT")
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, gitPushCmd, params, false, "")
	utils.CurlApp(t, *params)
	utils.Execute(t, cmd, params, false, "created initial release")
	utils.Execute(t, cmd, params, false, "listening on 5000...")

	utils.Execute(t, appsLogsLimitCmd, params, false, "")

	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
}

func appsOpenTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.CurlApp(t, *params)
	utils.CurlWithFail(t, fmt.Sprintf("http://%s.%s", "this-app-does-not-exist", params.Domain), true, "404 Not Found")
}

func appsRunTest(t *testing.T, params *utils.DeisTestConfig) {
	cmd := appsRunCmd
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.CheckList(t, cmd, params, "Hello, 世界", false)
	utils.Execute(t, "apps:run env", params, true, "GIT_SHA")
	// Fleet/systemd unit files have a limit of 2048 characters per line or else one encounters
	// problems parsing the unit.  To verify long log messages are truncated and do not crash
	// logspout (see https://github.com/deis/deis/issues/2046) we must issue a (relatively) short
	// command via `deis apps:run` that produces a LONG, but testable (predictable) log message we
	// can search for in the output of `deis logs`.
	//
	// The strategy for achieving this is to generate 1k random characters, then use that with a
	// command submitted via `deis apps:run` that will echo those 1k bytes 64x (on a single line).
	// Such a message is long enough to crash logspout if handled improperly and ALSO gives us a
	// large, distinct, and predictable string we can search for in the logs to assert success (and
	// assert that the message didn't crash logspout) WITHOUT ever needing to transmit such an
	// egregiously long command via `deis apps:run`.
	largeString := randomString(1024)
	utils.Execute(t, fmt.Sprintf("apps:run \"printf '%s%%.0s' {1..64}\"", largeString), params, false, largeString)
	// To assert the long message didn't crash logspout AND made it to the logger, we will search
	// the logs for a fragment of the long message-- specifically 2x the random string we generated.
	// This will help us ensure the actual log message made it through and not JUST the log message
	// that states the command being execured via `deis apps:run`.  We want to find the former, not
	// the latter because the latter is too short a message to have possibly crashed logspout if
	// mishandled.
	utils.Execute(t, "logs", params, false, strings.Repeat(largeString, 2))
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, cmd, params, true, "Not found")
}

func appsTransferTest(t *testing.T, params *utils.DeisTestConfig) {
	user := utils.GetGlobalConfig()
	user.UserName, user.Password = "app-transfer-test", "test"
	user.AppName = "transfer-test"
	user.NewOwner = params.UserName
	utils.Execute(t, authRegisterCmd, user, false, "")
	utils.Execute(t, authLoginCmd, user, false, "")
	utils.Execute(t, appsCreateCmdNoRemote, user, false, "")
	utils.Execute(t, appsTransferCmd, user, false, "")
	utils.Execute(t, appsInfoCmd, user, true, "403 FORBIDDEN")
	utils.Execute(t, authLoginCmd, params, false, "")
	params.AppName = user.AppName
	utils.CheckList(t, appsInfoCmd, params, params.UserName, false)
}
