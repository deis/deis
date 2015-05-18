package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/deis/deis/logger/syslogd"
)

var (
	logAddr         string
	logPort         int
	enablePublish   bool
	publishHost     string
	publishPath     string
	publishPort     string
	publishInterval int
	publishTTL      int
)

func init() {
	flag.StringVar(&logAddr, "log-addr", "0.0.0.0", "bind address for the logger")
	flag.IntVar(&logPort, "log-port", 514, "bind port for the logger")
	flag.StringVar(&syslogd.LogRoot, "log-root", "/data/logs", "log path to store logs")
	flag.BoolVar(&enablePublish, "enable-publish", false, "enable publishing to service discovery")
	flag.StringVar(&publishHost, "publish-host", getopt("HOST", "127.0.0.1"), "service discovery hostname")
	flag.IntVar(&publishInterval, "publish-interval", 10, "publish interval in seconds")
	flag.StringVar(&publishPath, "publish-path", getopt("ETCD_PATH", "/deis/logs"), "path to publish host/port information")
	flag.StringVar(&publishPort, "publish-port", getopt("ETCD_PORT", "4001"), "service discovery port")
	flag.IntVar(&publishTTL, "publish-ttl", publishInterval*2, "publish TTL in seconds")
}

func main() {
	flag.Parse()

	client := etcd.NewClient([]string{"http://" + publishHost + ":" + publishPort})

	// Wait for terminating signal
	exitChan := make(chan os.Signal, 2)
	cleanupChan := make(chan bool)
	signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT)

	go syslogd.Listen(exitChan, cleanupChan, fmt.Sprintf("%s:%d", logAddr, logPort))

	if enablePublish {
		go publishService(client, publishHost, publishPath, strconv.Itoa(logPort), uint64(time.Duration(publishTTL).Seconds()))
	}

	// Wait for the proper shutdown of the syslog server before exit
	<-cleanupChan
}

func publishService(client *etcd.Client, host string, etcdPath string, port string, ttl uint64) {
	for {
		setEtcd(client, etcdPath+"/host", host, ttl)
		setEtcd(client, etcdPath+"/port", port, ttl)
		time.Sleep(time.Duration(publishInterval))
	}
}

func setEtcd(client *etcd.Client, key, value string, ttl uint64) {
	_, err := client.Set(key, value, ttl)
	if err != nil && !strings.Contains(err.Error(), "Key already exists") {
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
