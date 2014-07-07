package verbose

import (
	"fmt"
	"github.com/deis/deis/tests/integration-utils"
	"testing"
)

func authSetup(t *testing.T) *itutils.UserDetails {
	cfg := itutils.GlobalSetup(t)
	ucfg := itutils.SetUser()
	fmt.Println("username :" + ucfg.UserName)
	fmt.Println("password :" + ucfg.Password)
	ucfg.HostName = cfg.HostName
	return ucfg
}

func authRegisterTest(t *testing.T, params *itutils.UserDetails) {
	cmd := itutils.GetCommand("auth", "register")
	itutils.Execute(t, cmd, params, false, "")
	itutils.Execute(t, cmd, params, true, "Registration failed")
}

func authLoginTest(t *testing.T, params *itutils.UserDetails) {
	cmd := itutils.GetCommand("auth", "login")
	itutils.Execute(t, cmd, params, false, "")
	params = authSetup(t)
	itutils.Execute(t, cmd, params, true, "200 OK")
}

func authLogoutTest(t *testing.T, params *itutils.UserDetails) {
	cmd := itutils.GetCommand("auth", "logout")
	itutils.Execute(t, cmd, params, false, "")

}

func authCancel() {
	fmt.Println("gexpect implementation")
}

func teardown(t *testing.T, params *itutils.UserDetails) {
	authLogoutTest(t, params)
}

func TestAuth(t *testing.T) {
	params := authSetup(t)
	authRegisterTest(t, params)
	authLogoutTest(t, params)
	authLoginTest(t, params)
	teardown(t, params)
}
