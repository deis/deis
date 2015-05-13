package main

import (
	"github.com/deis/deis/store/monitor/bindata"

	"github.com/deis/deis/pkg/boot"
	logger "github.com/deis/deis/pkg/log"
  "github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"
)

var (
	etcdPath = os.Getopt("ETCD_PATH", "/deis/store")
	log      = logger.New()
)

func init() {
	boot.RegisterComponent(new(ControllerBoot), "deis-component")
}

func main() {
	boot.Start(etcdPath, "-1", false)
}

type ControllerBoot struct{}

func (cb *ControllerBoot) MkdirsEtcd() []string {
	return []string{etcdPath}
}

func (cb *ControllerBoot) EtcdDefaults() map[string]string {
	numStores := os.Getopt("NUM_STORES", "3")
	pgNum := os.Getopt("PG_NUM", "128")
	keys := make(map[string]string)
	keys[etcdPath+"/size"] = numStores
	keys[etcdPath+"/minSize"] = "1"
	keys[etcdPath+"/pgNum"] = pgNum
	keys[etcdPath+"/delayStart"] = "15"
 	// We set this to the number of PGs before re-evaluating the PG count so users upgrading don't see the warning
	// Now, 12 pools * 64 pgs per pool = 768 PGs per OSD
	keys[etcdPath+"/maxPGsPerOSDWarning"] = "1536"
	return keys
}

func (cb *ControllerBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-store-monitor: starting...")

	setupParams := make(map[string]string)
	setupParams["ETCD_PATH"] = currentBoot.EtcdPath
	setupParams["ETCD"] = currentBoot.Host.String() + ":" + currentBoot.EtcdPort
	setupParams["HOST"] = currentBoot.Host.String()

	// TODO: this is required because PreBootScripts runs after confd.
	os.RunScript("monitor/bash/setup-monitor.bash", setupParams, bindata.Asset)
}

func (cb *ControllerBoot) PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	createParams := make(map[string]string)
	createParams["ETCD"] = currentBoot.Host.String() + ":" + currentBoot.EtcdPort
	return []*types.Script{
		&types.Script{Name: "monitor/bash/create-monitor.bash", Params: createParams, Content: bindata.Asset},
	}
}

func (cb *ControllerBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	hostname, _ := os.Hostname()
	cephMonCmd := "/usr/bin/ceph-mon -d -i " + hostname + " --public-addr " + currentBoot.Host.String() + ":6789"
	cmd, args := os.BuildCommandFromString(cephMonCmd)
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: cmd, Args: args}}
}

func (cb *ControllerBoot) WaitForPorts() []int {
	return []int{-1}
}

func (cb *ControllerBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (cb *ControllerBoot) PostBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-store-monitor: running...")
}

func (cb *ControllerBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	return []*types.Cron{}
}

func (cb *ControllerBoot) UseConfd() bool {
	return true
}
