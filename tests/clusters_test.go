// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/integration-utils"
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

func clustersUpdateTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("clusters", "update")
	// Regression test for https://github.com/deis/deis/pull/1283
	// Check that we didn't store the path of the key in the cluster.
	itutils.CheckList(t, params, cmd, "~/.ssh/", true)
}

func clustersDestroyTest(t *testing.T, params *itutils.DeisTestConfig) {
	cmd := itutils.GetCommand("clusters", "destroy")
	itutils.Execute(t, cmd, params, false, "")
}

func TestClusters(t *testing.T) {
	params := clustersSetup(t)
	clustersCreateTest(t, params)
	clustersListTest(t, params, false)
	clustersInfoTest(t, params)
	clustersUpdateTest(t, params)
	clustersDestroyTest(t, params)
	clustersListTest(t, params, true)

}
