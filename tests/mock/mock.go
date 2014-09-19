// Package mock provides mock objects and setup for Deis tests.

package mock

import (
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"
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
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "Starting")
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
