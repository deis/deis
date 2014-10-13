package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-etcd/etcd"

	"github.com/deis/deis/logger/syslogd"
)

const (
	timeout time.Duration = 10 * time.Second
	ttl     time.Duration = timeout * 2
)

func main() {
	host := getopt("HOST", "127.0.0.1")

	etcdPort := getopt("ETCD_PORT", "4001")
	etcdPath := getopt("ETCD_PATH", "/deis/logs")

	externalPort := getopt("EXTERNAL_PORT", "514")

	client := etcd.NewClient([]string{"http://" + host + ":" + etcdPort})

	// Wait for terminating signal
	exitChan := make(chan os.Signal, 2)
	cleanupChan := make(chan bool)
	signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT)

	go syslogd.Listen(exitChan, cleanupChan)

	go publishService(client, host, etcdPath, externalPort, uint64(ttl.Seconds()))

	// Wait for the proper shutdown of the syslog server before exit
	<-cleanupChan
}

func publishService(client *etcd.Client, host string, etcdPath string, externalPort string, ttl uint64) {
	for {
		setEtcd(client, etcdPath+"/host", host, ttl)
		setEtcd(client, etcdPath+"/port", externalPort, ttl)
		time.Sleep(timeout)
	}
}

func setEtcd(client *etcd.Client, key, value string, ttl uint64) {
	_, err := client.Set(key, value, ttl)
	if err != nil {
		log.Println(err)
	}
}

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}
