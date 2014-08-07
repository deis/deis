// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/integration-utils"
)

var (
	clustersCreateCmd  = "clusters:create {{.ClusterName}} {{.Domain}} --hosts={{.Hosts}} --auth={{.SSHKey}}"
	clustersListCmd    = "clusters:list"
	clustersUpdateCmd  = "clusters:update {{.ClusterName}} --domain={{.Domain}} --hosts={{.Hosts}} --auth=~/.ssh/{{.AuthKey}}"
	clustersInfoCmd    = "clusters:info {{.ClusterName}}"
	clustersDestroyCmd = "clusters:destroy {{.ClusterName}} --confirm={{.ClusterName}}"
)

func TestClusters(t *testing.T) {
	params := clustersSetup(t)
	clustersCreateTest(t, params)
	clustersListTest(t, params, false)
	clustersInfoTest(t, params)
	clustersUpdateTest(t, params)
	clustersDestroyTest(t, params)
	clustersListTest(t, params, true)
}

func clustersSetup(t *testing.T) *itutils.DeisTestConfig {
	cfg := itutils.GetGlobalConfig()
	cfg.ClusterName = "devtest"
	itutils.Execute(t, authLoginCmd, cfg, false, "")
	return cfg
}

func clustersCreateTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := clustersCreateCmd
	itutils.Execute(t, cmd, params, false, "")
	itutils.Execute(t, cmd, params, true, "Cluster with this Id already exists")
}

func clustersDestroyTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Execute(t, clustersDestroyCmd, params, false, "")
}

func clustersInfoTest(t *testing.T, params *itutils.DeisTestConfig) {
	itutils.Execute(t, clustersInfoCmd, params, false, "")
}

func clustersListTest(
	t *testing.T, params *itutils.DeisTestConfig, notflag bool) {
	itutils.CheckList(t, clustersListCmd, params, params.ClusterName, notflag)
}

func clustersUpdateTest(t *testing.T, params *itutils.DeisTestConfig) {
	// Regression test for https://github.com/deis/deis/pull/1283
	// Check that we didn't store the path of the key in the cluster.
	itutils.CheckList(t, clustersUpdateCmd, params, "~/.ssh/", true)
}
