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

// MesosBoot struct for mesos boot.
type MesosBoot struct{}

// MkdirsEtcd creates a directory in  etcd.
func (mb *MesosBoot) MkdirsEtcd() []string {
	return []string{etcdPath}
}

// EtcdDefaults returns default values for etcd.
func (mb *MesosBoot) EtcdDefaults() map[string]string {
	return map[string]string{}
}

// PreBootScripts runs preboot scripts.
func (mb *MesosBoot) PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

// PreBoot to log starting of marathon.
func (mb *MesosBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("mesos-slave: starting...")
}

// BootDaemons starts mesos-salve.
func (mb *MesosBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	args := gatherArgs(currentBoot.EtcdClient)
	args = append(args, "--ip="+currentBoot.Host.String())
	log.Infof("mesos slave args: %v", args)
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: "mesos-slave", Args: args}}
}

// WaitForPorts returns an array of ports.
func (mb *MesosBoot) WaitForPorts() []int {
	return []int{}
}

// PostBootScripts returns type script.
func (mb *MesosBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

// PostBoot returns type script.
func (mb *MesosBoot) PostBoot(currentBoot *types.CurrentBoot) {
	log.Info("mesos-slave: running...")
}

// ScheduleTasks returns a cron job.
func (mb *MesosBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	return []*types.Cron{}
}

// UseConfd uses confd.
func (mb *MesosBoot) UseConfd() (bool, bool) {
	return false, false
}

// PreShutdownScripts returns type script.
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
