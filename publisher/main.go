package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/fsouza/go-dockerclient"

	"github.com/deis/deis/publisher/server"
)

const (
	defaultRefreshTime time.Duration = 10 * time.Second
	defaultEtcdTTL     time.Duration = defaultRefreshTime * 2
	defaultHost                      = "127.0.0.1"
	defaultDockerHost                = "unix:///var/run/docker.sock"
	defaultEtcdHost                  = "127.0.0.1"
	defaultEtcdPort                  = "4001"
	defaultLogLevel                  = "error"
)

var (
	refreshDuration = flag.Duration("refresh-duration", defaultRefreshTime, "The time to wait between etcd refreshes.")
	etcdTTL         = flag.Duration("etcd-ttl", defaultEtcdTTL, "The TTL for all of the keys in etcd.")
	host            = flag.String("host", defaultHost, "The host where the publisher is running.")
	dockerHost      = flag.String("docker-host", defaultDockerHost, "The host where to find docker.")
	etcdPeersStr    = flag.String("etcd-peers", "", "The etcd peers.")
	etcdHost        = flag.String("etcd-host", defaultEtcdHost, "The etcd host.")
	etcdPort        = flag.String("etcd-port", defaultEtcdPort, "The etcd port.")
	logLevel        = flag.String("log-level", defaultLogLevel, "Acceptable values: error, debug")
)

func main() {

	flag.Parse()

	var etcdPeers []string
	if *etcdPeersStr == "" {
		etcdPeers = []string{"http://" + *etcdHost + ":" + *etcdPort}
	} else {
		etcdPeers = strings.Split(*etcdPeersStr, ",")
	}

	dockerClient, err := docker.NewClient(*dockerHost)
	if err != nil {
		log.Fatal(err)
	}
	etcdClient := etcd.NewClient(etcdPeers)

	server := server.New(dockerClient, etcdClient, *host, *logLevel)

	go server.Listen(*etcdTTL)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	for {
		go server.Poll(*etcdTTL)
		time.Sleep(*refreshDuration)
	}
}
