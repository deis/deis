package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/deis/deis/logger/syslogd"
)

var (
	enablePublish bool
	publishHost   string
	publishPath   string
	publishPort   string
	publishInterval int
	publishTTL      int
)

func init() {
	flag.BoolVar(&enablePublish, "publish", false, "enable publishing to service discovery")
	flag.IntVar(&publishInterval, "publish-interval", 10, "publish interval in seconds")
	flag.StringVar(&publishHost, "publish-host", getopt("HOST", "127.0.0.1"), "service discovery hostname")
	flag.StringVar(&publishPath, "publish-path", getopt("ETCD_PATH", "/deis/logs"), "path to publish host/port information")
	flag.StringVar(&publishPort, "publish-port", getopt("ETCD_PORT", "4001"), "service discovery port")
	flag.IntVar(&publishTTL, "publish-ttl", publishInterval*2, "publish TTL in seconds")
}

func main() {
	flag.Parse()

	externalPort := getopt("EXTERNAL_PORT", "514")

	client := etcd.NewClient([]string{"http://" + publishHost + ":" + publishPort})

	// Wait for terminating signal
	exitChan := make(chan os.Signal, 2)
	cleanupChan := make(chan bool)
	signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT)

	go syslogd.Listen(exitChan, cleanupChan)

	if enablePublish {
		go publishService(client, publishHost, publishPath, externalPort, uint64(time.Duration(publishTTL).Seconds()))
	}

	// Wait for the proper shutdown of the syslog server before exit
	<-cleanupChan
}

func publishService(client *etcd.Client, host string, etcdPath string, externalPort string, ttl uint64) {
	for {
		setEtcd(client, etcdPath+"/host", host, ttl)
		setEtcd(client, etcdPath+"/port", externalPort, ttl)
		time.Sleep(time.Duration(publishInterval))
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
