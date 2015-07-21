package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/deis/deis/mesos/bindata/zookeeper"
	"github.com/deis/deis/mesos/pkg/boot/zookeeper"
	"github.com/deis/deis/mesos/pkg/confd"
	"github.com/deis/deis/mesos/pkg/etcd"
	logger "github.com/deis/deis/mesos/pkg/log"
	oswrapper "github.com/deis/deis/mesos/pkg/os"
	"github.com/deis/deis/version"
)

var (
	etcdPath   = oswrapper.Getopt("ETCD_PATH", "/zookeeper/nodes")
	log        = logger.New()
	signalChan = make(chan os.Signal, 1)
)

func main() {
	host := oswrapper.Getopt("HOST", "127.0.0.1")
	etcdPort := oswrapper.Getopt("ETCD_PORT", "4001")
	etcdCtlPeers := oswrapper.Getopt("ETCD_PEERS", "127.0.0.1:"+etcdPort)
	etcdURL := etcd.GetHTTPEtcdUrls(host+":"+etcdPort, etcdCtlPeers)
	etcdClient := etcd.NewClient(etcdURL)

	etcd.Mkdir(etcdClient, etcdPath)

	log.Infof("boot version [%v]", version.Version)
	log.Info("zookeeper: starting...")

	zookeeper.CheckZkMappingInFleet(etcdPath, etcdClient, etcdURL)

	// we need to write the file /opt/zookeeper-data/data/myid with the id of this node
	os.MkdirAll("/opt/zookeeper-data/data", 0640)
	zkID := etcd.Get(etcdClient, etcdPath+"/"+host+"/id")
	ioutil.WriteFile("/opt/zookeeper-data/data/myid", []byte(zkID), 0640)

	zkServer := &zookeeper.ZkServer{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

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

	// wait for confd to run once and install initial templates
	confd.WaitForInitialConf(getConfdNodes(host, etcdCtlPeers, 4001), 10*time.Second)

	params := make(map[string]string)
	params["HOST"] = host
	if log.Level.String() == "debug" {
		params["DEBUG"] = "true"
	}

	err := oswrapper.RunScript("pkg/boot/zookeeper/bash/add-node.bash", params, bindata.Asset)
	if err != nil {
		log.Printf("command finished with error: %v", err)
	}

	if err := zkServer.Start(); err != nil {
		panic(err)
	}

	log.Info("zookeeper: running...")

	go func() {
		log.Debugf("starting pprof http server in port 6060")
		http.ListenAndServe("localhost:6060", nil)
	}()

	code := <-exitChan
	log.Debugf("execution terminated with exit code %v", code)

	log.Debugf("executing pre shutdown script")
	err = oswrapper.RunScript("pkg/boot/zookeeper/bash/remove-node.bash", params, bindata.Asset)
	if err != nil {
		log.Printf("command finished with error: %v", err)
	}

	log.Info("stopping zookeeper node")
	zkServer.Stop()
}

func getConfdNodes(host, etcdCtlPeers string, port int) []string {
	result := []string{host + ":" + strconv.Itoa(port)}

	if etcdCtlPeers != "127.0.0.1" {
		hosts := strings.Split(etcdCtlPeers, ",")
		result = []string{}
		for _, _host := range hosts {
			result = append(result, _host)
		}
	}

	return result
}
