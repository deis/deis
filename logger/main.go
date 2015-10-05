package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/deis/deis/logger/configurer"
	"github.com/deis/deis/logger/publisher"
	"github.com/deis/deis/logger/storage"
	"github.com/deis/deis/logger/syslogish"
	"github.com/deis/deis/logger/weblog"
)

var (
	// TODO: When semver permits us to do so, many of these flags should probably be phased out in
	// favor of just using environment variables.  Fewer avenues of configuring this component means
	// less confusion.
	logAddr       = flag.String("log-addr", "0.0.0.0", "bind address for the logger")
	logHost       = flag.String("log-host", getopt("HOST", "127.0.0.1"), "address of the host running logger")
	logPort       = flag.Int("log-port", 514, "bind port for the logger")
	enablePublish = flag.Bool("enable-publish", false, "enable publishing to service discovery")
	webAddr       = flag.String("web-addr", "0.0.0.0", "bind address for the web service")
	webPort       = flag.Int("web-port", 8088, "bind port for the web service")
	// Support legacy flag names even though the variable names have been changed for improved
	// clarity
	etcdHost        = flag.String("publish-host", getopt("HOST", "127.0.0.1"), "service discovery hostname")
	etcdPort        = flag.String("publish-port", getopt("ETCD_PORT", "4001"), "service discovery port")
	etcdPath        = flag.String("publish-path", getopt("ETCD_PATH", "/deis/logs"), "path to publish host/port information")
	configInterval  = flag.Int("config-interval", 10, "config interval in seconds")
	publishInterval = flag.Int("publish-interval", 10, "publish interval in seconds")
	publishTTL      int
)

func init() {
	flag.StringVar(&storage.LogRoot, "log-root", "/data/logs", "log path to store logs")
	// Support a legacy behavior that that allows default drain uri to be specified using a drain-uri
	// flag.
	flag.StringVar(&configurer.DefaultDrainURI, "drain-uri", "", "default drainURI, once set in etcd, this has no effect.")
	flag.Parse()
	// Set the default value for this AFTER the proper value of *publishInterval has been
	// established, since publishTTL should be twice the publishInterval by default.
	flag.IntVar(&publishTTL, "publish-ttl", *publishInterval*2, "publish TTL in seconds")
	// Now reparse flags in case the default publishTTL is overriden by a flag.
	flag.Parse()
}

func main() {
	syslogishServer, err := syslogish.NewServer(*logAddr, *logPort)
	if err != nil {
		log.Fatal("Error creating syslogish server", err)
	}
	weblogServer, err := weblog.NewServer(*webAddr, *webPort, syslogishServer)
	if err != nil {
		log.Fatal("Error creating weblog server", err)
	}
	etcdPortNum, err := strconv.Atoi(*etcdPort)
	if err != nil {
		log.Fatalf("Invalid port specified for etcd server.  '%s' is not an integer.", *etcdPort)
	}
	configurer, err := configurer.NewConfigurer(*etcdHost, etcdPortNum, *etcdPath, *configInterval,
		syslogishServer)
	if err != nil {
		log.Fatal("Error creating configurer", err)
	}

	configurer.Start()

	// Give configurer time to run once so we know syslogishServer is ready to rock
	time.Sleep(time.Duration(*configInterval+1) * time.Second)

	syslogishServer.Listen()
	weblogServer.Listen()

	if *enablePublish {
		publisher, err := publisher.NewPublisher(*etcdHost, etcdPortNum, *etcdPath, *publishInterval,
			publishTTL, *logHost, *logPort)
		if err != nil {
			log.Fatal("Error creating publisher", err)
		}
		publisher.Start()
	}

	log.Println("deis-logger running")

	// No cleanup is needed upon termination.  The signal to reopen log files (after hypothetical
	// logroation, for instance), if applicable, is the only signal we'll care about.  Our main loop
	// will just wait for that signal.
	reopen := make(chan os.Signal, 1)
	signal.Notify(reopen, syscall.SIGUSR1)

	for {
		<-reopen
		if err := syslogishServer.ReopenLogs(); err != nil {
			log.Fatal("Error reopening logs", err)
		}
	}
}

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}
