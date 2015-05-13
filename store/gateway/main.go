package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/deis/deis/store/gateway/bindata"

	"github.com/deis/deis/pkg/boot"
	logger "github.com/deis/deis/pkg/log"
	"github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"
)

const (
	servicePort = 8888
)

var (
	log          = logger.New()
	host         = os.Getopt("HOST", "127.0.0.1")
	etcdPath     = os.Getopt("ETCD_PATH", "/deis/store/gateway/hosts/"+host)
	externalPort = os.Getopt("EXTERNAL_PORT", strconv.Itoa(servicePort))
)

func init() {
	boot.RegisterComponent(new(ControllerBoot), "deis-component")
}

func main() {
	boot.Start(etcdPath, externalPort, false)
}

type ControllerBoot struct{}

func (cb *ControllerBoot) MkdirsEtcd() []string {
	return []string{etcdPath}
}

func (cb *ControllerBoot) EtcdDefaults() map[string]string {
	return map[string]string{}
}

func (cb *ControllerBoot) PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	setupParams := make(map[string]string)
	setupParams["ETCD_PATH"] = "/deis/store/gateway"
	setupParams["ETCD"] = currentBoot.Host.String() + ":" + currentBoot.EtcdPort
	setupParams["HOST"] = currentBoot.Host.String()
	return []*types.Script{
		&types.Script{Name: "gateway/bash/setup-gateway.bash", Params: setupParams, Content: bindata.Asset},
	}
}

func (cb *ControllerBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-store-gateway: starting...")
}

func (cb *ControllerBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	cmd, args := os.BuildCommandFromString("/usr/bin/radosgw -d -n client.radosgw.gateway")
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: cmd, Args: args}}
}

func (cb *ControllerBoot) WaitForPorts() []int {
	return []int{servicePort}
}

func (cb *ControllerBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (cb *ControllerBoot) PostBoot(currentBoot *types.CurrentBoot) {
	time.Sleep(5 * time.Second)
	log.Info("deis-store-gateway: radosgw running...")
}

func (cb *ControllerBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	params := make(map[string]string)
	params["HOST"] = currentBoot.Host.String()
	params["ETCD_PATH"] = "/deis/store/gateway"
	params["ETCD_TTL"] = fmt.Sprintf("%v", currentBoot.TTL.Seconds())
	params["EXTERNAL_PORT"] = strconv.Itoa(currentBoot.Port)
	params["ETCD"] = currentBoot.Host.String() + ":" + currentBoot.EtcdPort
	if log.Level.String() == "debug" {
		params["DEBUG"] = "true"
	}

	return []*types.Cron{
		&types.Cron{
			Frequency: "@every 5s",
			Code: func() {
				os.RunScript("gateway/bash/gateway-master.bash", params, bindata.Asset)
			},
		},
	}
}

func (cb *ControllerBoot) UseConfd() bool {
	return true
}
