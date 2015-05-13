//go:generate go-extpoints
package boot

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/deis/deis/pkg/boot/extpoints"

	"github.com/deis/deis/pkg/confd"
	"github.com/deis/deis/pkg/etcd"
	logger "github.com/deis/deis/pkg/log"
	netwrapper "github.com/deis/deis/pkg/net"
	oswrapper "github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"
	"github.com/robfig/cron"
	_ "net/http/pprof"
)

const (
	timeout time.Duration = 10 * time.Second
	ttl     time.Duration = timeout * 2
)

var (
	signalChan  = make(chan os.Signal, 1)
	log         = logger.New()
	bootProcess = extpoints.BootComponents
)

// RegisterComponent register an externsion to be used with this application
func RegisterComponent(component extpoints.BootComponent, name string) bool {
	return bootProcess.Register(component, name)
}

// Start initiate the boot process of the current component
// etcdPath is the base path used to publish the component in etcd
// externalPort is the base path used to publish the component in etcd
// useOneKeyIPPort indicates if we want to use just one key to publish the component
func Start(etcdPath string, externalPort string, useOneKeyIPPort bool) {
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
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
			default:
				exitChan <- 1
			}
		}
	}()

	// do the real work in a goroutine to be able to exit if
	// a signal is received during the boot process
	ep, _ := strconv.Atoi(externalPort)
	go start(etcdPath, ep, useOneKeyIPPort)

	code := <-exitChan
	log.Debugf("execution terminated with exit code %v", code)
	os.Exit(code)
}

func start(etcdPath string, externalPort int, useOneKeyIPPort bool) {
	component := bootProcess.Lookup("deis-component")
	if component == nil {
		log.Error("error loading deis extension...")
		signalChan <- syscall.SIGINT
	}

	log.Info("starting deis component...")

	host := oswrapper.Getopt("HOST", "127.0.0.1")
	etcdPort := oswrapper.Getopt("ETCD_PORT", "4001")
	etcdHostPort := host + ":" + etcdPort
	etcdClient := etcd.NewClient([]string{"http://" + etcdHostPort})

	currentBoot := &types.CurrentBoot{
		EtcdClient: etcdClient,
		EtcdPath:   etcdPath,
		EtcdPort:   etcdPort,
		Host:       net.ParseIP(host),
		Timeout:    timeout,
		TTL:        timeout * 2,
		Port:       externalPort,
	}

	if os.Getenv("DEBUG") != "" {
		go func() {
			listeningPort := netwrapper.RandomPort()
			log.Debugf("starting pprof http server in port %v", listeningPort)
			http.ListenAndServe("localhost:"+listeningPort, nil)
		}()
	}

	for _, key := range component.MkdirsEtcd() {
		etcd.Mkdir(etcdClient, key)
	}

	for key, value := range component.EtcdDefaults() {
		etcd.SetDefault(etcdClient, key, value)
	}

	component.PreBoot(currentBoot)

	if component.UseConfd() {
		// wait until etcd has discarded potentially stale values
		time.Sleep(timeout + 1)

		// wait for confd to run once and install initial templates
		confd.WaitForInitialConf(signalChan, etcdHostPort, timeout)
	}

	log.Debug("running pre boot scripts")
	preBootScripts := component.PreBootScripts(currentBoot)
	runAllScripts(signalChan, preBootScripts)

	if component.UseConfd() {
		// spawn confd in the background to update services based on etcd changes
		go confd.Launch(signalChan, etcdHostPort)
	}

	log.Debug("running boot daemons")
	servicesToStart := component.BootDaemons(currentBoot)
	for _, daemon := range servicesToStart {
		go oswrapper.RunProcessAsDaemon(signalChan, daemon.Command, daemon.Args)
	}

	portsToWaitFor := component.WaitForPorts()
	log.Debugf("waiting for a service in the port %v", portsToWaitFor)
	for _, portToWait := range portsToWaitFor {
		if portToWait > 0 {
			err := netwrapper.WaitForPort("tcp", "0.0.0.0", strconv.Itoa(portToWait), timeout)
			if err != nil {
				log.Printf("%v", err)
				signalChan <- syscall.SIGINT
			}
		}
	}

	// we only publish the service in etcd if the port if > 0
	if externalPort > 0 {
		log.Debug("starting periodic publication in etcd...")
		log.Debugf("etcd publication path %s, host %s and port %v", etcdPath, host, externalPort)
		// TODO: see another way to do this.
		// This is required because the router and store-gateway publish ip:port
		// in one key and not in different keys (host/port)
		if useOneKeyIPPort {
			go etcd.PublishServiceInOneKey(etcdClient, host, etcdPath, externalPort, uint64(ttl.Seconds()), timeout)
		} else {
			go etcd.PublishService(etcdClient, host, etcdPath, externalPort, uint64(ttl.Seconds()), timeout)
		}
	}

	// Wait for the first publication
	time.Sleep(timeout / 2)

	log.Printf("running post boot scripts")
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
			log.Printf("command finished with error: %v", err)
			signalChan <- syscall.SIGTERM
		}
	}
}
