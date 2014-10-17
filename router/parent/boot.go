package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/ActiveState/tail"
	"github.com/coreos/go-etcd/etcd"
)

const (
	timeout time.Duration = 10 * time.Second
	ttl     time.Duration = timeout * 2
)

func main() {
	host := getopt("HOST", "127.0.0.1")

	etcdPort := getopt("ETCD_PORT", "4001")

	etcdPath := getopt("ETCD_PATH", "/deis/router")

	hostEtcdPath := getopt("HOST_ETCD_PATH", "/deis/router/"+host)

	externalPort := getopt("EXTERNAL_PORT", "80")

	client := etcd.NewClient([]string{"http://" + host + ":" + etcdPort})

	// wait until etcd has discarded potentially stale values
	time.Sleep(timeout + 1)

	mkdirEtcd(client, "/deis/controller")
	mkdirEtcd(client, "/deis/services")
	mkdirEtcd(client, "/deis/domains")
	mkdirEtcd(client, "/deis/builder")

	setDefaultEtcd(client, etcdPath+"/gzip", "on")
	setDefaultEtcd(client, etcdPath+"/gzipHttpVersion", "1.1")
	setDefaultEtcd(client, etcdPath+"/gzipCompLevel", "5")
	setDefaultEtcd(client, etcdPath+"/gzipMinLength", "256")
	setDefaultEtcd(client, etcdPath+"/gzipProxied", "any")
	setDefaultEtcd(client, etcdPath+"/gzipVary", "on")
	setDefaultEtcd(client, etcdPath+"/gzipDisable", "\"msie6\"")
	setDefaultEtcd(client, etcdPath+"/gzipTypes", "\"application/atom+xml application/javascript application/json application/rss+xml application/vnd.ms-fontobject application/x-font-ttf application/x-web-app-manifest+json application/xhtml+xml application/xml font/opentype image/svg+xml image/x-icon text/css text/plain text/x-component\"")

	go publishService(client, host, hostEtcdPath, externalPort, uint64(ttl.Seconds()))

	log.Println("Starting Nginx...")
	go launchNginx()

	// tail the nginx logs and redirect to stdout and stderr
	go acessLog()
	go errorLog()

	waitForInitialConfd(host+":"+etcdPort, timeout)

	go launchConfd(host + ":" + etcdPort)

	exitChan := make(chan os.Signal, 2)
	signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT)
	go cleanOnExit(exitChan)
	<-exitChan
}

func cleanOnExit(exitChan chan os.Signal) {
	for _ = range exitChan {
		tail.Cleanup()
	}
}

func waitForInitialConfd(etcd string, timeout time.Duration) {
	for {
		cmd := exec.Command("confd", "-onetime", "-node", etcd, "-config-file", "/app/confd.toml")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err == nil {
			break
		}

		time.Sleep(timeout)
	}
}

func launchConfd(etcd string) {
	cmd := exec.Command("confd", "-node", etcd, "-config-file", "/app/confd.toml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Printf("confd terminated by error: %v", err)
	}
}

func launchNginx() {
	cmd := exec.Command("/opt/nginx/sbin/nginx", "-c", "/opt/nginx/conf/nginx.conf")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Printf("Nginx terminated by error: %v", err)
	}

	// Wait until the nginx is available
	for {
		_, err := net.DialTimeout("tcp", "127.0.0.1:80", timeout)
		if err == nil {
			log.Println("deis-router running...")
			break
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Nginx terminated by error: %v", err)
	} else {
		log.Println("Reloading nginx (change in configuration)...")
	}
}

func acessLog() {
	t, _ := tail.TailFile("/opt/nginx/logs/access.log", tail.Config{Follow: true})

	for line := range t.Lines {
		fmt.Fprintln(os.Stdout, line.Text)
	}
}

func errorLog() {
	t, _ := tail.TailFile("/opt/nginx/logs/error.log", tail.Config{Follow: true})

	for line := range t.Lines {
		fmt.Fprintln(os.Stderr, line.Text)
	}
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

func setDefaultEtcd(client *etcd.Client, key, value string) {
	_, err := client.Set(key, value, 0)
	if err != nil {
		log.Println(err)
	}
}

func mkdirEtcd(client *etcd.Client, path string) {
	_, err := client.CreateDir(path, 0)
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
