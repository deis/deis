package verbose

import (
	_ "fmt"
	"github.com/deis/deis/tests/integration-utils"
	"testing"
)

func clustersSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	cfg.ClusterName = "devtest"
	cmd := itutils.GetCommand("auth", "login")
	itutils.Execute(t, cmd, cfg, false, "")
	return cfg
}

func clustersCreateTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("clusters", "create")
	itutils.Execute(t, cmd, params, false, "")
	itutils.Execute(t, cmd, params, true, "Cluster with this Id already exists")
}

func clustersListTest(t *testing.T, params *itutils.DeisTestConfig, notflag bool) {
	cmd := itutils.GetCommand("clusters", "list")
	itutils.CheckList(t, params, cmd, params.ClusterName, notflag)
}

func clustersInfoTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("clusters", "info")
	itutils.Execute(t, cmd, params, false, "")
}

//Tets #1283

func clustersUpdateTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("clusters", "update")
	itutils.CheckList(t, params, cmd, "~/.ssh/"+params.AuthKey, true)
}

func clustersDestroyTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("clusters", "destroy")
	itutils.Execute(t, cmd, params, false, "")
}

func TestKeys(t *testing.T) {
	params := clustersSetup(t)
	clustersCreateTest(t, params)
	clustersListTest(t, params, false)
	clustersInfoTest(t, params)
	clustersUpdateTest(t, params)
	clustersDestroyTest(t, params)
	clustersListTest(t, params, true)

}
