package main

import (
	"strings"

	"github.com/deis/deis/mesos/pkg/boot"
	"github.com/deis/deis/mesos/pkg/etcd"
	logger "github.com/deis/deis/mesos/pkg/log"
	"github.com/deis/deis/mesos/pkg/os"
	"github.com/deis/deis/mesos/pkg/types"
)

const (
	mesosPort = 5050
)

var (
	etcdPath = os.Getopt("ETCD_PATH", "/mesos/master")
	log      = logger.New()
)

func init() {
	boot.RegisterComponent(new(MesosBoot), "boot")
}

func main() {
	boot.Start(etcdPath, mesosPort)
}

type MesosBoot struct{}

func (mb *MesosBoot) MkdirsEtcd() []string {
	return []string{etcdPath}
}

func (mb *MesosBoot) EtcdDefaults() map[string]string {
	return map[string]string{}
}

func (mb *MesosBoot) PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (mb *MesosBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("mesos-master: starting...")
}

func (mb *MesosBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	args := gatherArgs(currentBoot.EtcdClient)
	args = append(args, "--ip="+currentBoot.Host.String())
	args = append(args, "--hostname="+currentBoot.Host.String())
	log.Infof("mesos master args: %v", args)
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: "mesos-master", Args: args}}
}

func (mb *MesosBoot) WaitForPorts() []int {
	return []int{}
}

func (mb *MesosBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (mb *MesosBoot) PostBoot(currentBoot *types.CurrentBoot) {
	log.Info("mesos-master: running...")
}

func (mb *MesosBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	return []*types.Cron{}
}

func (mb *MesosBoot) UseConfd() (bool, bool) {
	return false, false
}

func (mb *MesosBoot) PreShutdownScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func gatherArgs(c *etcd.Client) []string {
	var args []string

	nodes := etcd.GetList(c, "/zookeeper/nodes")
	var hosts []string
	for _, node := range nodes {
		hosts = append(hosts, node+":3888")
	}
	zkHosts := strings.Join(hosts, ",")
	args = append(args, "--zk=zk://"+zkHosts+"/mesos")
	args = append(args, "--quorum=1")
	args = append(args, "--work_dir=/tmp/mesos-master")

	return args
}
