package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

const (
	timeout         time.Duration = 10 * time.Second
	ttl             time.Duration = timeout * 2
	redisWait       time.Duration = 5 * time.Second
	redisConf       string        = "/app/redis.conf"
	ectdKeyNotFound int           = 100
	defaultMemory   string        = "50mb"
)

func main() {
	host := getopt("HOST", "127.0.0.1")

	etcdPort := getopt("ETCD_PORT", "4001")
	etcdPath := getopt("ETCD_PATH", "/deis/cache")

	externalPort := getopt("EXTERNAL_PORT", "6379")

	client := etcd.NewClient([]string{"http://" + host + ":" + etcdPort})

	var maxmemory string
	result, err := client.Get("/deis/cache/maxmemory", false, false)
	if err != nil {
		if e, ok := err.(*etcd.EtcdError); ok && e.ErrorCode == ectdKeyNotFound {
			maxmemory = defaultMemory
		} else {
			log.Fatalln(err)
		}
	} else {
		maxmemory = result.Node.Key
	}
	replaceMaxmemoryInConfig(maxmemory)

	go launchRedis()

	go publishService(client, host, etcdPath, externalPort, uint64(ttl.Seconds()))

	// Wait for terminating signal
	exitChan := make(chan os.Signal, 2)
	signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT)
	<-exitChan
}

func replaceMaxmemoryInConfig(maxmemory string) {
	input, err := ioutil.ReadFile(redisConf)
	if err != nil {
		log.Fatalln(err)
	}
	output := strings.Replace(string(input), "# maxmemory <bytes>", "maxmemory "+maxmemory, 1)
	err = ioutil.WriteFile(redisConf, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func launchRedis() {
	cmd := exec.Command("redis-server", redisConf)
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
