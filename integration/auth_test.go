package verbose

import (
	"fmt"
	"github.com/deis/deis/tests/integration-utils"
	"testing"
)

func authSetup(t *testing.T) *AuthData {
	_ := itutils.GlobalSetup(t)
	ucfg := itutils.SetUser()
	fmt.Println("username :" + ucfg.UserName)
	fmt.Println("password :" + ucfg.Password)
	return &ucfg
}

func authRegisterTest(t *testing.T, params *AuthData) {
	cmd := itutils.GetCommand("auth", "register")
	itutils.Execute(t, cmd, params, false)
	itutils.Execute(t, cmd, params, true)
}

func authLoginTest(t *testing.T, params *AuthData) {
	cmd := itutils.GetCommand("auth", "login")
	itutils.Execute(t, cmd, params, false)
	params = authSetup(t)
	itutils.Execute(t, cmd, params, true)
}

func authLogoutTest(t *testing.T, params *AuthData) {
	cmd := itutils.GetCommand("auth", "logout")
	itutils.Execute(t, cmd, params, false)

}

func authCancel() {
	fmt.Println("coming soon")
}

func TestAuth(t *testing.T) {
	params := authSetup(t)
	authRegisterTest(t, params)
	authLogoutTest(t, params)
	authLoginTest(t, params)

}
