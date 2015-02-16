package main

import (
	"log"
	"os"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/fsouza/go-dockerclient"

	"github.com/deis/deis/publisher/server"
)

const (
	timeout time.Duration = 10 * time.Second
	etcdTTL time.Duration = timeout * 2
)

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}

func main() {
	endpoint := getopt("DOCKER_HOST", "unix:///var/run/docker.sock")
	etcdHost := getopt("ETCD_HOST", "127.0.0.1")

	client, err := docker.NewClient(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	etcdClient := etcd.NewClient([]string{"http://" + etcdHost + ":4001"})

	server := &server.Server{DockerClient: client, EtcdClient: etcdClient}

	go server.Listen(etcdTTL)

	for {
		go server.Poll(etcdTTL)
		time.Sleep(timeout)
	}
}
