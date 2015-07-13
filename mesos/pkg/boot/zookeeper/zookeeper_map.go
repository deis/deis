package zookeeper

import (
	"strconv"

	"github.com/deis/deis/mesos/pkg/etcd"
	"github.com/deis/deis/mesos/pkg/fleet"
	logger "github.com/deis/deis/mesos/pkg/log"
)

const (
	etcdLock = "/zookeeper/setupLock"
)

var (
	log = logger.New()
)

// CheckZkMappingInFleet verifies if there is a mapping for each node in
// the CoreOS cluster using the metadata zookeeper=true to filter wich
// nodes zookeeper should run
func CheckZkMappingInFleet(etcdPath string, etcdClient *etcd.Client, etcdURL []string) {
	// check if the nodes with the required role already have the an id.
	// If not get fleet nodes with the required role and preassing the
	// ids for every node in the cluster
	err := etcd.AcquireLock(etcdClient, etcdLock, 10)
	if err != nil {
		panic(err)
	}

	zkNodes := etcd.GetList(etcdClient, etcdPath)
	log.Debugf("zookeeper nodes %v", zkNodes)

	machines, err := getMachines(etcdURL)
	if err != nil {
		panic(err)
	}
	log.Debugf("machines %v", machines)

	if len(machines) == 0 {
		log.Warning("")
		log.Warning("there is no machine using metadata zookeeper=true in the cluster to run zookeeper")
		log.Warning("we will create the mapping with for all the nodes")
		log.Warning("")
		machines = fleet.GetNodesInCluster(etcdURL)
	}

	if len(zkNodes) == 0 {
		log.Debug("initializing zookeeper cluster")
		for index, newZkNode := range machines {
			log.Debug("adding node %v to zookeeper cluster", newZkNode)
			etcd.Set(etcdClient, etcdPath+"/"+newZkNode+"/id", strconv.Itoa(index+1), 0)
		}
	} else {
		// we check if some machine in the fleet cluster with the
		// required role is not initialized (no zookeeper node id).
		machinesNotInitialized := difference(machines, zkNodes)
		if len(machinesNotInitialized) > 0 {
			nextNodeID := getNextNodeID(etcdPath, etcdClient, zkNodes)
			for _, zkNode := range machinesNotInitialized {
				etcd.Set(etcdClient, etcdPath+"/"+zkNode+"/id", strconv.Itoa(nextNodeID), 0)
				nextNodeID++
			}
		}
	}

	// release the etcd lock
	etcd.ReleaseLock(etcdClient)
}

// getMachines return the list of machines that can run zookeeper or an empty list
func getMachines(etcdURL []string) ([]string, error) {
	metadata, err := fleet.ParseMetadata("zookeeper=true")
	if err != nil {
		panic(err)
	}

	return fleet.GetNodesWithMetadata(etcdURL, metadata)
}

// getNextNodeID returns the next id to use as zookeeper node index
func getNextNodeID(etcdPath string, etcdClient *etcd.Client, nodes []string) int {
	result := 0
	for _, node := range nodes {
		id := etcd.Get(etcdClient, etcdPath+"/"+node+"/id")
		numericID, err := strconv.Atoi(id)
		if id != "" && err == nil && numericID > result {
			result = numericID
		}
	}

	return result + 1
}

// difference get the elements present in the first slice and not in
// the second one returning those elemenets in a new string slice.
func difference(slice1 []string, slice2 []string) []string {
	diffStr := []string{}
	m := map[string]int{}

	for _, s1Val := range slice1 {
		m[s1Val] = 1
	}
	for _, s2Val := range slice2 {
		m[s2Val] = m[s2Val] + 1
	}

	for mKey, mVal := range m {
		if mVal == 1 {
			diffStr = append(diffStr, mKey)
		}
	}

	return diffStr
}
