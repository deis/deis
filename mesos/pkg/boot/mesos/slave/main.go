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
	mesosPort = 5051
)

var (
	etcdPath = os.Getopt("ETCD_PATH", "/mesos/slave")
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
	log.Info("mesos-slave: starting...")
}

func (mb *MesosBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	args := gatherArgs(currentBoot.EtcdClient)
	args = append(args, "--ip="+currentBoot.Host.String())
	log.Infof("mesos slave args: %v", args)
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: "mesos-slave", Args: args}}
}

func (mb *MesosBoot) WaitForPorts() []int {
	return []int{}
}

func (mb *MesosBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (mb *MesosBoot) PostBoot(currentBoot *types.CurrentBoot) {
	log.Info("mesos-slave: running...")
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
	args = append(args, "--master=zk://"+zkHosts+"/mesos")
	args = append(args, "--containerizers=docker,mesos")

	return args
}
