package main

import (
	"strings"

	"github.com/deis/deis/mesos/pkg/boot/mesos/marathon/bindata"

	"github.com/deis/deis/mesos/pkg/boot"
	"github.com/deis/deis/mesos/pkg/etcd"
	logger "github.com/deis/deis/mesos/pkg/log"
	"github.com/deis/deis/mesos/pkg/os"
	"github.com/deis/deis/mesos/pkg/types"
)

const (
	mesosPort = 8180
)

var (
	etcdPath = os.Getopt("ETCD_PATH", "/mesos/marathon")
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
	params := make(map[string]string)
	params["HOST"] = currentBoot.Host.String()
	err := os.RunScript("pkg/boot/mesos/marathon/bash/update-hosts-file.bash", params, bindata.Asset)
	if err != nil {
		log.Printf("command finished with error: %v", err)
	}

	return []*types.Script{}
}

func (mb *MesosBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("mesos-marathon: starting...")
}

func (mb *MesosBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	args := gatherArgs(currentBoot.EtcdClient)
	log.Infof("mesos marathon args: %v", args)
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: "/marathon/bin/start", Args: args}}
}

func (mb *MesosBoot) WaitForPorts() []int {
	return []int{}
}

func (mb *MesosBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (mb *MesosBoot) PostBoot(currentBoot *types.CurrentBoot) {
	log.Info("mesos-marathon: running...")
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
	args = append(args, "--master", "zk://"+zkHosts+"/mesos")
	args = append(args, "--zk", "zk://"+zkHosts+"/marathon")
	// 20min task launch timeout for large docker image pulls
	args = append(args, "--task_launch_timeout", "1200000")
	args = append(args, "--ha")
	args = append(args, "--http_port", "8180")

	return args
}
