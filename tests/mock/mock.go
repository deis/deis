// Package mock provides mock objects and setup for Deis tests.

package mock

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"

	"github.com/docker/docker/api/client"
)

// RunMockDatabase starts a mock postgresql database for testing.
func RunMockDatabase(t *testing.T, tag string, etcdPort string, dbPort string) {
	var err error
	cli, stdout, stdoutPipe := dockercli.NewClient()
	done := make(chan bool, 1)
	dbImage := "deis/test-postgresql:latest"
	ipaddr := utils.HostAddress()
	done <- true
	go func() {
		<-done
		err = dockercli.RunContainer(cli,
			"--name", "deis-test-database-"+tag,
			"--rm",
			"-p", dbPort+":5432",
			"-e", "EXTERNAL_PORT="+dbPort,
			"-e", "HOST="+ipaddr,
			"-e", "USER=deis",
			"-e", "DB=deis",
			"-e", "PASS=deis",
			dbImage)
	}()
	time.Sleep(1000 * time.Millisecond)
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "Initialization complete.")
	setkeys := []string{
		"/deis/database/user",
		"/deis/database/password",
		"/deis/database/name",
	}
	setdir := []string{}
	dbhandler := etcdutils.InitEtcd(setdir, setkeys, etcdPort)
	etcdutils.PublishEtcd(t, dbhandler)
	etcdutils.SetEtcd(t,
		[]string{"/deis/database/host", "/deis/database/port", "/deis/database/engine"},
		[]string{ipaddr, dbPort, "postgresql_psycopg2"}, dbhandler.C)
	if err != nil {
		t.Fatal(err)
	}
}

// RunMockCeph runs a set of containers used to mock a Ceph storage cluster
func RunMockCeph(t *testing.T, name string, cli *client.DockerCli, etcdPort string) {

	etcdutils.SetSingle(t, "/deis/store/hosts/"+utils.HostAddress(), utils.HostAddress(), etcdPort)

	monitorName := name + "-monitor"
	RunMockCephMonitor(t, monitorName, "6379", etcdPort)

	daemonName := name + "-daemon"
	RunMockCephDaemon(t, daemonName, "6800", etcdPort)

	gatewayName := name + "-gateway"
	RunMockCephGateway(t, gatewayName, utils.RandomPort(), etcdPort)
}

// RunMockCephMonitor runs a Ceph Monitor agent
func RunMockCephMonitor(t *testing.T, name string, port string, etcdPort string) {
	var err error
	cli, stdout, stdoutPipe := dockercli.NewClient()
	cephImage := "deis/store-monitor:" + utils.BuildTag()
	ipaddr := utils.HostAddress()
	cephAddr := ipaddr + ":" + port
	fmt.Printf("--- Running deis/mock-ceph-monitor at %s\n", cephAddr)
	done2 := make(chan bool, 1)
	go func() {
		done2 <- true
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-p", port,
			"-e", "HOST="+ipaddr,
			"-e", "ETCD_PORT="+etcdPort,
			"-e", "NUM_STORES=1",
			"--net=host",
			cephImage)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "monmap e1: 1 mons at")
	if err != nil {
		t.Fatal(err)
	}
}

// RunMockCephDaemon sets up a single Ceph OSD
func RunMockCephDaemon(t *testing.T, name string, port string, etcdPort string) {
	var err error
	cli, stdout, stdoutPipe := dockercli.NewClient()
	cephImage := "deis/store-daemon:" + utils.BuildTag()
	ipaddr := utils.HostAddress()
	cephAddr := ipaddr + ":" + port
	fmt.Printf("--- Running deis/mock-ceph-daemon at %s\n", cephAddr)
	done := make(chan bool, 1)
	go func() {
		done <- true
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-p", port,
			"-e", "HOST="+ipaddr,
			"-e", "ETCD_PORT="+etcdPort,
			"--net=host",
			cephImage)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "journal close /var/lib/ceph/osd/ceph-0/journal")
	if err != nil {
		t.Fatal(err)
	}
}

// RunMockCephGateway starts a mock S3 endpoint used for component testing
func RunMockCephGateway(t *testing.T, name string, port string, etcdPort string) {
	var err error
	cli, stdout, stdoutPipe := dockercli.NewClient()
	cephImage := "deis/store-gateway:" + utils.BuildTag()
	ipaddr := utils.HostAddress()
	cephAddr := ipaddr + ":" + port
	fmt.Printf("--- Running deis/mock-ceph-gateway at %s\n", cephAddr)
	done2 := make(chan bool, 1)
	go func() {
		done2 <- true
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"-h", "deis-store-gateway",
			"--rm",
			"-p", port+":"+"8888",
			"-e", "ETCD_PORT="+etcdPort,
			"-e", "HOST="+ipaddr,
			"-e", "EXTERNAL_PORT="+port,
			cephImage)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "deis-store-gateway running...")
	if err != nil {
		t.Fatal(err)
	}
}
