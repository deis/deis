package main

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/deis/deis/pkg/boot"
	logger "github.com/deis/deis/pkg/log"
	oswrapper "github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"

	"github.com/ActiveState/tail"
)

const (
	servicePort = 80
)

var (
	hostEtcdPath   = oswrapper.Getopt("HOST_ETCD_PATH", "/deis/router/hosts")
	externalPort   = oswrapper.Getopt("EXTERNAL_PORT", strconv.Itoa(servicePort))
	etcdPath       = oswrapper.Getopt("ETCD_PATH", "/deis/router")
	gitLogFile     = "/opt/nginx/logs/git.log"
	nginxAccessLog = "/opt/nginx/logs/access.log"
	nginxErrorLog  = "/opt/nginx/logs/error.log"
	log            = logger.New()
)

func init() {
	boot.RegisterComponent(new(RouterBoot), "deis-component")
}

func main() {
	boot.Start(hostEtcdPath, externalPort, true)
}

type RouterBoot struct{}

func (rb *RouterBoot) MkdirsEtcd() []string {
	return []string{
		etcdPath,
		"/deis/controller",
		"/deis/services",
		"/deis/domains",
		"/deis/builder",
		"/deis/router/hosts",
		"/deis/certs",
	}
}

func (rb *RouterBoot) EtcdDefaults() map[string]string {
	keys := make(map[string]string)
	keys[etcdPath+"/gzip"] = "on"
	return keys
}

func (rb *RouterBoot) PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (rb *RouterBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-router: starting...")
	go tailFile(nginxAccessLog)
	go tailFile(nginxErrorLog)
	go tailFile(gitLogFile)
}

func (rb *RouterBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	nginxCommand := "/opt/nginx/sbin/nginx -c /opt/nginx/conf/nginx.conf"
	cmd, args := oswrapper.BuildCommandFromString(nginxCommand)
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: cmd, Args: args}}
}

func (rb *RouterBoot) WaitForPorts() []int {
	return []int{servicePort}
}

func (rb *RouterBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	return []*types.Script{}
}

func (rb *RouterBoot) PostBoot(currentBoot *types.CurrentBoot) {
	log.Info("deis-router: nginx running...")
}

func (rb *RouterBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	return []*types.Cron{
		&types.Cron{
			Frequency: "@every 30s",
			Code: func() {
				cmd := exec.Command("/bin/generate-certs")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			},
		},
	}
}

func (rb *RouterBoot) UseConfd() bool {
	return true
}

func tailFile(path string) {
	mkfifo(path)
	t, _ := tail.TailFile(path, tail.Config{Follow: true})

	for line := range t.Lines {
		log.Info(line.Text)
	}
}

func mkfifo(path string) {
	os.Remove(path)
	if err := syscall.Mkfifo(path, syscall.S_IFIFO|0666); err != nil {
		log.Fatalf("%v", err)
	}
}
