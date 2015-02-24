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
	var err error
	cli, stdout, stdoutPipe := dockercli.NewClient()
	cephImage := "deis/mock-store:latest"
	ipaddr := utils.HostAddress()
	fmt.Printf("--- Running deis/mock-store at %s\n", ipaddr)
	done2 := make(chan bool, 1)
	go func() {
		done2 <- true
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-e", "HOST="+ipaddr,
			"-e", "ETCD_PORT="+etcdPort,
			"--net=host",
			cephImage)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "deis-store-gateway running...")
	if err != nil {
		t.Fatal(err)
	}
}
