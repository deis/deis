package main

import (
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-etcd/etcd"

	"github.com/deis/deis/boot/builder"
	"github.com/deis/deis/boot/commons"
	"github.com/deis/deis/boot/logger"
)

const (
	timeout time.Duration = 10 * time.Second
	ttl     time.Duration = timeout * 2
	tcpWait time.Duration = timeout / 2
)

func main() {
	logger.Log.Info("Starting deis-builder...")

	host := commons.Getopt("HOST", "127.0.0.1")

	etcdPort := commons.Getopt("ETCD_PORT", "4001")
	etcdPath := commons.Getopt("ETCD_PATH", "/deis/builder")

	externalPort := commons.Getopt("EXTERNAL_PORT", "2223")

	storageDriver := commons.Getopt("STORAGE_DRIVER", "btrfs")

	client := etcd.NewClient([]string{"http://" + host + ":" + etcdPort})

	commons.MkdirEtcd(client, etcdPath)
	commons.MkdirEtcd(client, etcdPath+"/users")

	// wait until etcd has discarded potentially stale values
	time.Sleep(timeout + 1)

	etcdHostPort := host + ":" + etcdPort

	// check for stored configuration in deis-store
	builder.CheckSSHKeysInStore(client)

	// wait for confd to run once and install initial templates
	commons.WaitForInitialConfd(etcdHostPort, timeout)

	// spawn confd in the background to update services based on etcd changes
	commons.LaunchConfd(etcdHostPort)

	// remove any pre-existing docker.sock
	// spawn a docker daemon to run builds
	os.Remove("/var/run/docker.sock")

	go launchDocker(storageDriver)

	// wait for docker
	waitForDocker()

	// HACK: load progrium/cedarish tarball for faster boot times
	// see https://github.com/deis/deis/issues/1027
	checkCedarish()

	logger.Log.Println("building slugbuilder and slugrunner...")
	buildImage("deis/slugbuilder", "/app/slugbuilder/")
	buildImage("deis/slugrunner", "/app/slugrunner/")

	// start an SSH daemon to process `git push` requests
	go launchSshd()

	go commons.PublishService(client, host, etcdPath, externalPort, uint64(ttl.Seconds()), timeout)

	// Wait for terminating signal
	exitChan := make(chan os.Signal, 2)
	signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT)
	<-exitChan
}

func waitForDocker() {
	logger.Log.Debug("waiting for docker daemon to be available...")
	for {
		cmd := exec.Command("docker", "info")
		if err := cmd.Run(); err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func launchSshd() {
	logger.Log.Debug("starting ssh server...")
	cmd := exec.Command("/usr/sbin/sshd", "-D", "-e")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		logger.Log.Fatalf("Error starting SSH server: %v", err)
	}

	// Wait until the ssh server is available
	for {
		if _, err := net.DialTimeout("tcp", "127.0.0.1:22", tcpWait); err == nil {
			break
		}
	}

	time.Sleep(tcpWait)
	logger.Log.Info("deis-builder running...")

	err = cmd.Wait()
	logger.Log.Infof("sshd finished by error: %v", err)
}

func launchDocker(storageDriver string) {
	logger.Log.Debug("starting docker daemon...")
	cmd := exec.Command("docker", "-d", "--storage-driver="+storageDriver, "--bip=172.19.42.1/16")
	err := cmd.Start()
	if err != nil {
		logger.Log.Fatalf("Error starting docker daemon: %v", err)
	}

	logger.Log.Debug("docker daemon started...")
	err = cmd.Wait()
	logger.Log.Infof("sshd finished by error: %v", err)
}

func checkCedarish() {
	logger.Log.Debug("checking for cedarish...")
	cmd := exec.Command("docker", "history", "progrium/cedarish")
	if err := cmd.Run(); err != nil {
		logger.Log.Println("loading cedarish...")
		cmd := exec.Command("docker", "load", "-i", "/progrium_cedarish.tar")
		if err := cmd.Run(); err != nil {
			logger.Log.Fatal(err)
		}
	} else {
		logger.Log.Println("cedarish already loaded")
	}
}

func buildImage(tagName string, directory string) {
	logger.Log.Debugf("building image %s...", tagName)
	cmd := exec.Command("docker", "build", "-t", tagName, directory)
	err := cmd.Run()
	if err != nil {
		logger.Log.Fatal(err)
	}
}
