package main

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/deis/deis/pkg/boot"
	"github.com/deis/deis/pkg/etcd"
	logger "github.com/deis/deis/pkg/log"
	"github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"
)

const (
	redisConf     = "/app/redis.conf"
	defaultMemory = "50mb"
	servicePort   = 6379
)

var (
	etcdPath     = os.Getopt("ETCD_PATH", "/deis/cache")
	externalPort = os.Getopt("EXTERNAL_PORT", strconv.Itoa(servicePort))
	log          = logger.New()
	memory       string
)

func init() {
	boot.RegisterComponent(new(CacheBoot), "deis-component")
}

func main() {
	boot.Start(etcdPath, externalPort, false)
}

type CacheBoot struct{}

func (cb *CacheBoot) MkdirsEtcd() []string {
	return []string{etcdPath}
}

func (cb *CacheBoot) EtcdDefaults() map[string]string {
	return map[string]string{}
}

func (cb *CacheBoot) PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (cb *CacheBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-cache: starting...")
	memory := etcd.Get(currentBoot.EtcdClient, "/deis/cache/maxmemory")
	if memory == "" {
		memory = defaultMemory
	}
	replaceMaxmemoryInConfig(memory)
}

func (cb *CacheBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	cmd, args := os.BuildCommandFromString("redis-server " + redisConf)
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: cmd, Args: args}}
}

func (cb *CacheBoot) WaitForPorts() []int {
	return []int{servicePort}
}

func (cb *CacheBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (cb *CacheBoot) PostBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-cache: redis is running...")
}

func (cb *CacheBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	return []*types.Cron{}
}

func (cb *CacheBoot) UseConfd() bool {
	return false
}

func replaceMaxmemoryInConfig(maxmemory string) {
	input, err := ioutil.ReadFile(redisConf)
	if err != nil {
		log.Fatalln(err)
	}
	output := strings.Replace(string(input), "# maxmemory <bytes>", "maxmemory "+maxmemory, 1)
	err = ioutil.WriteFile(redisConf, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
