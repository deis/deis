package verbose

import (
	"fmt"
	"github.com/ThomasRooney/gexpect"
	"github.com/deis/deis/tests/integration-utils"
	"github.com/deis/deis/tests/utils"
	"testing"
)

func authSetup(t *testing.T) *itutils.DeisTestConfig {
	user := itutils.GetGlobalConfig()
	user.UserName, user.Password = utils.GetUserDetails()
	fmt.Println("username :" + user.UserName)
	fmt.Println("password :" + user.Password)
	return user
}

func authRegisterTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("auth", "register")
	itutils.Execute(t, cmd, params, false, "")
	itutils.Execute(t, cmd, params, true, "Registration failed")
}

func authLoginTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("auth", "login")
	itutils.Execute(t, cmd, params, false, "")
	params = authSetup(t)
	itutils.Execute(t, cmd, params, true, "200 OK")
}

func authLogoutTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("auth", "logout")
	itutils.Execute(t, cmd, params, false, "")
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
