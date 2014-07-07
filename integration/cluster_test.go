package verbose

import (
	_ "fmt"
	"github.com/deis/deis/tests/integration-utils"
	"testing"
)

func clustersSetup(t *testing.T) *ClusterDetails {
	cfg := itutils.GlobalSetup(t)
	cscfg := itutils.ClusterDetails{
		cfg.ClusterName,
		cfg.Hosts,
		"172.17.8.100",
		cfg.AuthKey,
		cfg.HostName,
	}
	cmd := itutils.GetCommand("keys", "add")
	itutils.Execute(t, cmd, cfg, false, "")
	return &cscfg
}

func clustersCreateTest(t *testing.T, params *ClusterDetails) {
	cmd := itutils.GetCommand("clusters", "create")
	itutils.Execute(t, cmd, params, false, "")
	itutils.Execute(t, cmd, params, true, "Cluster with this Id already exists")
}

func clustersListTest(t *testing.T, params *ClusterDetails) {
	cmd := itutils.GetCommand("clusters", "list")
	itutils.Execute(t, cmd, params, false, "")
}

func clustersInfoTest(t *testing.T, params *ClusterDetails) {
	cmd := itutils.GetCommand("clusters", "info")
	itutils.Execute(t, cmd, params, false, "")
	params.ClusterName = "kin"
	itutils.Execute(t, cmd, params, true, "Not found")
	params.ClusterName = "dev"
}

func clustersUpdateTest(t *testing.T, params *ClusterDetails) {
	cmd := itutils.GetCommand("clusters", "update")
	itutils.Execute(t, cmd, params, false, "")
	params.ClusterName = "kin"
	itutils.Execute(t, cmd, params, true, "Not found")
	params.ClusterName = "dev"
}

func clustersDestroyTest(t *testing.T, params *ClusterDetails) {
	cmd := itutils.GetCommand("clusters", "destroy")
	itutils.Execute(t, cmd, params, false, "")
	itutils.Execute(t, cmd, params, true, "Not found")
}

func TestKeys(t *testing.T) {
	params := clustersSetup(t)
	clustersCreateTest(t, params)
	clustersListTest(t, params)
	clustersInfoTest(t, params)
	clustersUpdateTest(t, params)
	//clustersDestroyTest(t, params)

}
