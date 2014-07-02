package verbose

import (
	"fmt"
	"github.com/deis/deis/tests/integration-utils"
	"testing"
)

func keysSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GlobalSetup(t)
	return cfg
}

func keysAddTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("keys", "add")
	itutils.Execute(t, cmd, params, false,"")
	itutils.Execute(t, cmd, params, true,"Uploading deis to Deis...400 BAD REQUEST")
}

func keysListTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("keys", "list")
	itutils.Execute(t, cmd, params, false,"")
}

func keysRemoveTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("keys", "remove")
	itutils.Execute(t, cmd, params, false,"")
  itutils.Execute(t, cmd, params, true,"Not found")
}

func authCancel() {
	fmt.Println("coming soon")
}

func TestKeys(t *testing.T) {
	params := keysSetup(t)
	keysAddTest(t, params)
	keysListTest(t, params)
	keysRemoveTest(t, params)
}
