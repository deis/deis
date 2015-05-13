package main

import (
	"io/ioutil"
	"os/exec"
	"strconv"
	"time"

	"github.com/deis/deis/builder/bindata"

	"github.com/deis/deis/pkg/boot"
	logger "github.com/deis/deis/pkg/log"
	"github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"
)

const (
	servicePort = 22
)

var (
	etcdPath     = os.Getopt("ETCD_PATH", "/deis/builder")
	externalPort = os.Getopt("EXTERNAL_PORT", strconv.Itoa(servicePort))
	log          = logger.New()
)

func init() {
	boot.RegisterComponent(new(BuilderBoot), "deis-component")
}

func main() {
	boot.Start(etcdPath, externalPort, false)
}

type BuilderBoot struct{}

func (bb *BuilderBoot) MkdirsEtcd() []string {
	return []string{
		etcdPath,
		etcdPath + "/users",
	}
}

func (bb *BuilderBoot) EtcdDefaults() map[string]string {
	return map[string]string{}
}

func (bb *BuilderBoot) PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{
		&types.Script{Name: "bash/copy-apparmor.bash", Content: bindata.Asset},
		&types.Script{Name: "bash/check-overlay.bash", Content: bindata.Asset},
		&types.Script{Name: "bash/build-slugbuilder-slugrunner.bash", Content: bindata.Asset},
	}
}

func (bb *BuilderBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-builder: starting...")
}

func (bb *BuilderBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	driverOverride := readDockerEnvFile()
	dockerArgs := []string{
		"--daemon",
		"--bip=172.19.42.1/16",
		"--insecure-registry=10.0.0.0/8",
		"--insecure-registry=172.16.0.0/12",
		"--insecure-registry=192.168.0.0/16",
		"--insecure-registry=100.64.0.0/10",
	}

	if driverOverride != "" {
		log.Debugf("custom docker env [%v]", driverOverride)
		dockerArgs = append(dockerArgs, driverOverride)
	}

	log.Debugf("starting docker daemon...")
	sshCmd, sshArgs := os.BuildCommandFromString("/usr/sbin/sshd -D -e")
	return []*types.ServiceDaemon{
		&types.ServiceDaemon{Command: "/usr/bin/docker", Args: dockerArgs},
		&types.ServiceDaemon{Command: sshCmd, Args: sshArgs},
	}
}

func (bb *BuilderBoot) WaitForPorts() []int {
	return []int{servicePort}
}

func (bb *BuilderBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (bb *BuilderBoot) PostBoot(currentBoot *types.CurrentBoot) {
	waitForDocker()
	log.Info("deis-builder: running...")
}

func (bb *BuilderBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	return []*types.Cron{}
}

func (bb *BuilderBoot) UseConfd() bool {
	return true
}

func waitForDocker() {
	log.Debug("waiting for docker daemon to be available...")
	for {
		cmd := exec.Command("docker", "info")
		if err := cmd.Run(); err == nil {
			break
		}

		time.Sleep(5 * time.Second)
	}
}

func readDockerEnvFile() string {
	envFile, err := ioutil.ReadFile("/etc/docker.env")
	if err != nil {
		log.Debugf("%v", err)
		return ""
	}
	return string(envFile)
}
