package fleet

import (
	"net/http"
	"time"

	"github.com/coreos/fleet/etcd"
	"github.com/coreos/fleet/registry"
	logger "github.com/deis/deis/mesos/pkg/log"
)

var log = logger.New()

// GetNodesWithMetadata returns the ip address of the nodes with all the specified roles
func GetNodesWithMetadata(url []string, metadata map[string][]string) ([]string, error) {
	etcdClient, err := etcd.NewClient(url, &http.Transport{}, time.Second)
	if err != nil {
		log.Debugf("error creating new fleet etcd client: %v", err)
		return nil, err
	}

	fleetClient := registry.NewEtcdRegistry(etcdClient, "/_coreos.com/fleet/")
	machines, err := fleetClient.Machines()
	if err != nil {
		log.Debugf("error creating new fleet etcd client: %v", err)
		return nil, err
	}

	var machineList []string
	for _, m := range machines {
		if hasMetadata(m, metadata) {
			machineList = append(machineList, m.PublicIP)
		}
	}

	return machineList, nil
}

// GetNodesInCluster return the list of ip address of all the nodes
// running in the cluster currently active (fleetctl list-machines)
func GetNodesInCluster(url []string) []string {
	etcdClient, err := etcd.NewClient(url, &http.Transport{}, time.Second)
	if err != nil {
		log.Debugf("error creating new fleet etcd client: %v", err)
		return []string{}
	}

	fleetClient := registry.NewEtcdRegistry(etcdClient, "/_coreos.com/fleet/")
	machines, err := fleetClient.Machines()
	if err != nil {
		log.Debugf("error creating new fleet etcd client: %v", err)
		return []string{}
	}

	var machineList []string
	for _, m := range machines {
		machineList = append(machineList, m.PublicIP)
	}

	return machineList
}
