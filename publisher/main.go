package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/deis/deis/publisher/server"
)

var (
	bindAddr   string
	interval   time.Duration
	publishTTL time.Duration
	host       string
	dockerAddr string
	etcdAddr   string
	log        = logrus.New()
	logLevel   string
)

func init() {
	flag.StringVar(&bindAddr, "bind-addr", "localhost:6060", "address to listen for incoming HTTP requests")
	flag.DurationVar(&interval, "interval", 10*time.Second, "backend polling interval")
	flag.DurationVar(&publishTTL, "publish-ttl", interval*2, "backend TTL when publishing keys")
	flag.StringVar(&host, "host", "127.0.0.1", "host address of the machine")
	flag.StringVar(&dockerAddr, "docker-addr", "unix:///var/run/docker.sock", "address to a docker API")
	flag.StringVar(&etcdAddr, "etcd-addr", "http://127.0.0.1:4001", "address to the etcd host")
	flag.StringVar(&logLevel, "log-level", "error", "level which publisher should log messages")
}

func main() {
	flag.Parse()

	signalChan := make(chan os.Signal, 2)
	stopChan := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info("booting publisher...")
	setLogLevel()

	// run a profiler
	go func() {
		if err := http.ListenAndServe(bindAddr, nil); err != nil {
			log.Warningf("failed to run profiler (%s)", err)
		} else {
			log.Infof("profiler listening on %s", bindAddr)
		}
	}()

	server, err := server.New(dockerAddr, etcdAddr, host, publishTTL)
	if err != nil {
		log.Fatalf("failed to start (%s)", err)
	}

	// run Publish() once at boot to publish existing containers
	go server.Publish()
	go server.Listen(stopChan)

	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			// we need to re-publish containers every interval so the TTL doesn't expire
			go server.Publish()
		case <-signalChan:
			log.Info("shutting down publisher...")
			stopChan <- true
			break
		}
	}
}

func setLogLevel() {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("failed to parse log level '%s' (%s)", logLevel, err)
	}
	logrus.SetLevel(lvl)
}
