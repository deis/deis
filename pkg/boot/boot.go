package boot

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deis/deis/pkg/confd"
	"github.com/deis/deis/pkg/etcd"
	Log "github.com/deis/deis/pkg/log"
	. "github.com/deis/deis/pkg/net"
	. "github.com/deis/deis/pkg/os"
	_ "net/http/pprof"
)

const (
	timeout time.Duration = 10 * time.Second
)

var (
	signalChan = make(chan os.Signal, 2)
	log        = Log.New()
)

// New contructor that indicates the etcd base path and
// the port that the component will expose
func New(etcdPath, port string) *boot {
	log.Info("starting deis component...")

	host := Getopt("HOST", "127.0.0.1")
	etcdPort := Getopt("ETCD_PORT", "4001")

	etcdHostPort := host + ":" + etcdPort

	etcdClient := etcd.NewClient([]string{"http://" + etcdHostPort})

	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	if os.Getenv("DEBUG") != "" {
		go func() {
			http.ListenAndServe("localhost:6060", nil)
		}()
	}

	return &boot{
		Etcd:     etcdClient,
		EtcdPath: etcdPath,
		EtcdPort: etcdPort,
		Host:     net.ParseIP(host),
		Timeout:  timeout,
		TTL:      timeout * 2,
		Port:     port,
	}
}

// StartConfd initiates the boot process waiting for the correct initialization
// of the required values for the confd template and launching confd as daemon
func (boot *boot) StartConfd() {
	// wait until etcd has discarded potentially stale values
	time.Sleep(boot.Timeout + 1)

	// wait for confd to run once and install initial templates
	confd.WaitForInitialConf(signalChan, boot.Host.String()+":"+boot.EtcdPort, boot.Timeout)

	// spawn confd in the background to update services based on etcd changes
	go confd.Launch(signalChan, boot.Host.String()+":"+boot.EtcdPort)
}

// Publish publish information about the relevant process running in the boot
// process in etcd using specified path and port/s
func (boot *boot) Publish(port ...string) {
	portToPublish := boot.Port
	// If we specify a custom port we use that one
	if len(port) != 0 {
		portToPublish = port[1]
	}
	log.Debug("starting periodic publication in etcd...")
	log.Debugf("etcd publication path %s, host %s and port %s", boot.EtcdPath, boot.Host, portToPublish)
	go etcd.PublishService(boot.Etcd, boot.Host.String(), boot.EtcdPath, portToPublish, uint64(boot.TTL.Seconds()), boot.Timeout)
}

// RunProcessAsDaemon pkg/os RunProcessAsDaemon wrapper
func (boot *boot) RunProcessAsDaemon(command string, args []string) {
	go RunProcessAsDaemon(signalChan, command, args)
}

// RunScript pkg/os RunScript wrapper
func (boot *boot) RunScript(script string, params map[string]string, loader func(string) ([]byte, error)) {
	RunScript(signalChan, script, params, loader)
}

// WaitForLocalConnection wait until the port/ports exposed are opened
// If no port is specified we use the defined in the constructor
func (boot *boot) WaitForLocalConnection(ports ...string) {
	if len(ports) == 0 {
		log.Debugf("waiting for a service in the port %v", boot.Port)
		WaitForPort("tcp", "127.0.0.1", boot.Port, boot.Timeout)
	} else {
		// we need to wait for a port different than the default or more than one
		log.Debugf("waiting for the services in the port/s %v", ports)
		for _, port := range ports {
			WaitForPort("tcp", "127.0.0.1", port, boot.Timeout)
		}
	}
}

// Wait wait until a SIGTERM or SIGINT signal is received
func (boot *boot) Wait() {
	<-signalChan
}
