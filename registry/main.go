package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/deis/deis/registry/bindata"

	"github.com/deis/deis/pkg/boot"
	logger "github.com/deis/deis/pkg/log"
	"github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"
)

const (
	servicePort = 5000
)

var (
	host         = os.Getopt("HOST", "127.0.0.1")
	etcdPath     = os.Getopt("ETCD_PATH", "/deis/registry/hosts/"+host)
	externalPort = os.Getopt("EXTERNAL_PORT", strconv.Itoa(servicePort))
	log          = logger.New()
)

func init() {
	boot.RegisterComponent(new(RegistryBoot), "deis-component")
}

func main() {
	boot.Start(etcdPath, externalPort, false)
}

type RegistryBoot struct{}

func (rb *RegistryBoot) MkdirsEtcd() []string {
	return []string{
		etcdPath,
		"/deis/store/gateway",
	}
}

func (rb *RegistryBoot) EtcdDefaults() map[string]string {
	bucketName := os.Getopt("BUCKET_NAME", "registry")
	keys := make(map[string]string)
	keys["/deis/registry/protocol"] = "http"
	keys["/deis/registry/bucketName"] = bucketName
	return keys
}

func (rb *RegistryBoot) PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	bucketName := os.Getopt("BUCKET_NAME", "registry")
	params := make(map[string]string)
	params["BUCKET_NAME"] = bucketName
	return []*types.Script{
		&types.Script{Name: "bash/create-bucket.bash", Params: params, Content: bindata.Asset},
	}
}

func (rb *RegistryBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-registry: starting...")
}

func (rb *RegistryBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	cmd, args := os.BuildCommandFromString("sudo -E -u registry docker-registry")
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: cmd, Args: args}}
}

func (rb *RegistryBoot) WaitForPorts() []int {
	return []int{servicePort}
}

func (rb *RegistryBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (rb *RegistryBoot) PostBoot(currentBoot *types.CurrentBoot) {
	time.Sleep(5 * time.Second)
	log.Info("deis-registry: docker-registry is running...")
}

func (rb *RegistryBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	params := make(map[string]string)
	params["HOST"] = currentBoot.Host.String()
	params["ETCD_PATH"] = "/deis/registry"
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
				os.RunScript("bash/registry-master.bash", params, bindata.Asset)
			},
		},
	}
}

func (rb *RegistryBoot) UseConfd() bool {
	return true
}
