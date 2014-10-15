package main

import (
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

const (
	timeout   time.Duration = 10 * time.Second
	ttl       time.Duration = timeout * 2
	redisWait time.Duration = 5 * time.Second
)

func main() {
	host := getopt("HOST", "127.0.0.1")

	etcdPort := getopt("ETCD_PORT", "4001")
	etcdPath := getopt("ETCD_PATH", "/deis/cache")

	externalPort := getopt("EXTERNAL_PORT", "6379")

	client := etcd.NewClient([]string{"http://" + host + ":" + etcdPort})

	go launchRedis()

	go publishService(client, host, etcdPath, externalPort, uint64(ttl.Seconds()))

	// Wait for terminating signal
	exitChan := make(chan os.Signal, 2)
	signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT)
	<-exitChan
}

func launchRedis() {
	cmd := exec.Command("/app/bin/redis-server", "/app/redis.conf")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if err != nil {
		log.Printf("Error starting Redis: %v", err)
		os.Exit(1)
	}

	// Wait until the redis server is available
	for {
		_, err := net.DialTimeout("tcp", "127.0.0.1:6379", redisWait)
		if err == nil {
			log.Println("deis-cache running...")
			break
		}
	}

	err = cmd.Wait()
	log.Printf("Redis finished by error: %v", err)
}

func publishService(client *etcd.Client, host string, etcdPath string,
	externalPort string, ttl uint64) {

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
