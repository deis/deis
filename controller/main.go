package main

import (
	"strconv"

	"github.com/deis/deis/controller/bindata"

	"github.com/deis/deis/pkg/boot"
	logger "github.com/deis/deis/pkg/log"
	"github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"
)

const (
	servicePort = 8000
)

var (
	etcdPath     = os.Getopt("ETCD_PATH", "/deis/controller")
	externalPort = os.Getopt("EXTERNAL_PORT", strconv.Itoa(servicePort))
	log          = logger.New()
)

func init() {
	boot.RegisterComponent(new(ControllerBoot), "deis-component")
}

func main() {
	boot.Start(etcdPath, externalPort, false)
}

type ControllerBoot struct{}

func (cb *ControllerBoot) MkdirsEtcd() []string {
	return []string{
		etcdPath,
		"/deis/services",
		"/deis/domains",
		"/deis/platform",
		"/deis/scheduler",
	}
}

func (cb *ControllerBoot) EtcdDefaults() map[string]string {
	protocol := os.Getopt("DEIS_PROTOCOL", "http")
	sk, _ := os.Random(64)
	secretKey := os.Getopt("DEIS_SECRET_KEY", sk)
	bk, _ := os.Random(64)
	builderKey := os.Getopt("DEIS_BUILDER_KEY", bk)
	keys := make(map[string]string)
	keys[etcdPath+"/protocol"] = protocol
	keys[etcdPath+"/secretKey"] = secretKey
	keys[etcdPath+"/builderKey"] = builderKey
	keys[etcdPath+"/registrationEnabled"] = "1"
	keys[etcdPath+"/webEnabled"] = "0"
	keys[etcdPath+"/unitHostname"] = "default"
	return keys
}

func (cb *ControllerBoot) PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{
		&types.Script{Name: "bash/migrate.bash", Content: bindata.Asset},
	}
}

func (cb *ControllerBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-controller: starting...")
}

func (cb *ControllerBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	cmd, args := os.BuildCommandFromString("sudo -E -u deis gunicorn -c /app/deis/gconf.py deis.wsgi")
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: cmd, Args: args}}
}

func (cb *ControllerBoot) WaitForPorts() []int {
	return []int{servicePort}
}

func (cb *ControllerBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (cb *ControllerBoot) PostBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-controller: running...")
}

func (cb *ControllerBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	return []*types.Cron{}
}

func (cb *ControllerBoot) UseConfd() bool {
	return true
}
