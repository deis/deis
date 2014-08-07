// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/integration-utils"
	"github.com/deis/deis/tests/utils"
)

var (
	authLoginCmd    = "auth:login http://deis.{{.Domain}} --username={{.UserName}} --password={{.Password}}"
	authLogoutCmd   = "auth:logout"
	authRegisterCmd = "auth:register http://deis.{{.Domain}} --username={{.UserName}} --password={{.Password}} --email={{.Email}}"
)

func TestAuth(t *testing.T) {
	params := authSetup(t)
	authRegisterTest(t, params)
	authLogoutTest(t, params)
	authLoginTest(t, params)
	authCancel(t, params)
}

func authSetup(t *testing.T) *itutils.DeisTestConfig {
	user := itutils.GetGlobalConfig()
	user.UserName, user.Password = utils.NewID(), utils.NewID()
	return user
}

func authCancel(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.AuthCancel(t, params)
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

func authRegisterTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := authRegisterCmd
	itutils.Execute(t, cmd, params, false, "")
	itutils.Execute(t, cmd, params, true, "Registration failed")
}
