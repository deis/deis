//go:generate go-extpoints

package boot

import (
	"net/http"
	_ "net/http/pprof" //pprof is used for profiling servers
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/deis/deis/mesos/pkg/boot/extpoints"
	"github.com/deis/deis/mesos/pkg/confd"
	"github.com/deis/deis/mesos/pkg/etcd"
	logger "github.com/deis/deis/mesos/pkg/log"
	"github.com/deis/deis/mesos/pkg/net"
	oswrapper "github.com/deis/deis/mesos/pkg/os"
	"github.com/deis/deis/mesos/pkg/types"
	"github.com/deis/deis/version"
	"github.com/robfig/cron"
)

const (
	timeout time.Duration = 10 * time.Second
	ttl     time.Duration = timeout * 2
)

var (
	signalChan  = make(chan os.Signal, 1)
	log         = logger.New()
	bootProcess = extpoints.BootComponents
	component   extpoints.BootComponent
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// RegisterComponent register an externsion to be used with this application
func RegisterComponent(component extpoints.BootComponent, name string) bool {
	return bootProcess.Register(component, name)
}

// Start initiate the boot process of the current component
// etcdPath is the base path used to publish the component in etcd
// externalPort is the base path used to publish the component in etcd
func Start(etcdPath string, externalPort int) {
	log.Infof("boot version [%v]", version.Version)

	go func() {
		log.Debugf("starting pprof http server in port 6060")
		http.ListenAndServe("localhost:6060", nil)
	}()

	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt,
	)

	// Wait for a signal and exit
	exitChan := make(chan int)
	go func() {
		for {
			s := <-signalChan
			log.Debugf("Signal received: %v", s)
			switch s {
			case syscall.SIGTERM:
				exitChan <- 0
			case syscall.SIGQUIT:
				exitChan <- 0
			case syscall.SIGKILL:
				exitChan <- 1
			default:
				exitChan <- 1
			}
		}
	}()

	component = bootProcess.Lookup("boot")
	if component == nil {
		log.Error("error loading boot extension...")
		signalChan <- syscall.SIGINT
	}

	host := oswrapper.Getopt("HOST", "127.0.0.1")
	etcdPort, _ := strconv.Atoi(oswrapper.Getopt("ETCD_PORT", "4001"))
	etcdPeers := oswrapper.Getopt("ETCD_PEERS", "127.0.0.1:"+strconv.Itoa(etcdPort))
	etcdClient := etcd.NewClient(etcd.GetHTTPEtcdUrls(host+":"+strconv.Itoa(etcdPort), etcdPeers))

	etcdURL := etcd.GetHTTPEtcdUrls(host+":"+strconv.Itoa(etcdPort), etcdPeers)

	currentBoot := &types.CurrentBoot{
		ConfdNodes: getConfdNodes(host+":"+strconv.Itoa(etcdPort), etcdPeers),
		EtcdClient: etcdClient,
		EtcdPath:   etcdPath,
		EtcdPort:   etcdPort,
		EtcdPeers:  etcdPeers,
		EtcdURL:    etcdURL,
		Host:       net.ParseIP(host),
		Timeout:    timeout,
		TTL:        timeout * 2,
		Port:       externalPort,
	}

	// do the real work in a goroutine to be able to exit if
	// a signal is received during the boot process
	go start(currentBoot)

	code := <-exitChan

	// pre shutdown tasks
	log.Debugf("executing pre shutdown scripts")
	preShutdownScripts := component.PreShutdownScripts(currentBoot)
	runAllScripts(signalChan, preShutdownScripts)

	log.Debugf("execution terminated with exit code %v", code)
	os.Exit(code)
}

func start(currentBoot *types.CurrentBoot) {
	log.Info("starting component...")

	log.Debug("creating required etcd directories")
	for _, key := range component.MkdirsEtcd() {
		etcd.Mkdir(currentBoot.EtcdClient, key)
	}

	log.Debug("setting default etcd values")
	for key, value := range component.EtcdDefaults() {
		etcd.SetDefault(currentBoot.EtcdClient, key, value)
	}

	// component.PreBoot(currentBoot)

	initial, daemon := component.UseConfd()
	if initial {
		// wait for confd to run once and install initial templates
		log.Debug("waiting for initial confd configuration")
		confd.WaitForInitialConf(currentBoot.ConfdNodes, currentBoot.Timeout)
	}

	log.Debug("running preboot code")
	component.PreBoot(currentBoot)

	log.Debug("running pre boot scripts")
	preBootScripts := component.PreBootScripts(currentBoot)
	runAllScripts(signalChan, preBootScripts)

	if daemon {
		// spawn confd in the background to update services based on etcd changes
		log.Debug("launching confd")
		go confd.Launch(signalChan, currentBoot.ConfdNodes)
	}

	log.Debug("running boot daemons")
	servicesToStart := component.BootDaemons(currentBoot)
	for _, daemon := range servicesToStart {
		go oswrapper.RunProcessAsDaemon(signalChan, daemon.Command, daemon.Args)
	}

	// if the returned ips contains the value contained in $HOST it means
	// that we are running docker with --net=host
	ipToListen := "0.0.0.0"
	netIfaces := net.GetNetworkInterfaces()
	for _, iface := range netIfaces {
		if strings.Index(iface.IP, currentBoot.Host.String()) > -1 {
			ipToListen = currentBoot.Host.String()
			break
		}
	}

	portsToWaitFor := component.WaitForPorts()
	log.Debugf("waiting for a service in the port %v in ip %v", portsToWaitFor, ipToListen)
	for _, portToWait := range portsToWaitFor {
		if portToWait > 0 {
			err := net.WaitForPort("tcp", ipToListen, portToWait, timeout)
			if err != nil {
				log.Errorf("error waiting for port %v using ip %v: %v", portToWait, ipToListen, err)
				signalChan <- syscall.SIGINT
			}
		}
	}

	time.Sleep(60 * time.Second)

	// we only publish the service in etcd if the port if > 0
	if currentBoot.Port > 0 {
		log.Debug("starting periodic publication in etcd...")
		log.Debugf("etcd publication path %s, host %s and port %v", currentBoot.EtcdPath, currentBoot.Host, currentBoot.Port)
		go etcd.PublishService(currentBoot.EtcdClient, currentBoot.EtcdPath+"/"+currentBoot.Host.String(), currentBoot.Host.String(), currentBoot.Port, uint64(ttl.Seconds()), timeout)

		// Wait for the first publication
		time.Sleep(timeout / 2)
	}

	log.Debug("running post boot scripts")
	postBootScripts := component.PostBootScripts(currentBoot)
	runAllScripts(signalChan, postBootScripts)

	log.Debug("checking for cron tasks...")
	crons := component.ScheduleTasks(currentBoot)
	_cron := cron.New()
	for _, cronTask := range crons {
		_cron.AddFunc(cronTask.Frequency, cronTask.Code)
	}
	_cron.Start()

	component.PostBoot(currentBoot)
}

func getConfdNodes(host, etcdCtlPeers string) []string {
	if etcdCtlPeers != "127.0.0.1:4001" {
		hosts := strings.Split(etcdCtlPeers, ",")
		result := []string{}
		for _, _host := range hosts {
			result = append(result, _host)
		}
		return result
	}

	return []string{host}
}

func runAllScripts(signalChan chan os.Signal, scripts []*types.Script) {
	for _, script := range scripts {
		if script.Params == nil {
			script.Params = map[string]string{}
		}
		// add HOME variable to avoid warning from ceph commands
		script.Params["HOME"] = "/tmp"
		if log.Level.String() == "debug" {
			script.Params["DEBUG"] = "true"
		}
		err := oswrapper.RunScript(script.Name, script.Params, script.Content)
		if err != nil {
			log.Errorf("script %v execution finished with error: %v", script.Name, err)
			signalChan <- syscall.SIGTERM
		}
	}
}
