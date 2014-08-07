// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/integration-utils"
	"github.com/deis/deis/tests/utils"
)

var (
	authLoginCmd = `auth:login http://deis.{{.Domain}} \
--username={{.UserName}} --password={{.Password}}`
	authLogoutCmd   = "auth:logout"
	authRegisterCmd = `auth:register http://deis.{{.Domain}} \
--username={{.UserName}} --password={{.Password}} --email={{.Email}}`
)

func authSetup(t *testing.T) *itutils.DeisTestConfig {
	user := itutils.GetGlobalConfig()
	user.UserName, user.Password = utils.GetUserDetails()
	return user
}

func authRegisterTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := authRegisterCmd
	itutils.Execute(t, cmd, params, false, "")
	itutils.Execute(t, cmd, params, true, "Registration failed")
}

func authLoginTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := authLoginCmd
	itutils.Execute(t, cmd, params, false, "")
	params = authSetup(t)
	itutils.Execute(t, cmd, params, true, "200 OK")
}

func authLogoutTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Execute(t, authLogoutCmd, params, false, "")
}

func authCancel(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.AuthCancel(t, params)
}

func teardown(t *testing.T, params *itutils.DeisTestConfig) {
	authLogoutTest(t, params)
}

func TestAuth(t *testing.T) {
	params := authSetup(t)
	authRegisterTest(t, params)
	authLogoutTest(t, params)
	authLoginTest(t, params)
	authCancel(t, params)
}
