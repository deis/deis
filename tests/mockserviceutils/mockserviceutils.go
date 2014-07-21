package mockserviceutils

import (
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"
)

// RunMockDatabase starts a mock postgresql database for testing.
func RunMockDatabase(t *testing.T, uid string, etcdPort string, dbPort string) {
	var err error
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dbImage := "paintedfox/postgresql:latest"
	ipaddr := utils.GetHostIPAddress()
	done <- true
	go func() {
		<-done
		err = dockercliutils.RunContainer(cli,
			"--name", "deis-test-database-"+uid,
			"--rm",
			"-p", dbPort+":5432",
			"-e", "PUBLISH="+dbPort,
			"-e", "HOST="+ipaddr,
			"-e", "USER=deis",
			"-e", "DB=deis",
			"-e", "PASS=deis",
			dbImage)
	}()
	time.Sleep(1000 * time.Millisecond)
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "Starting")
	setkeys := []string{
		"/deis/database/user",
		"/deis/database/password",
		"/deis/database/name",
	}
	setdir := []string{}
	dbhandler := etcdutils.InitetcdValues(setdir, setkeys, etcdPort)
	etcdutils.Publishvalues(t, dbhandler)
	etcdutils.SetEtcdValues(t,
		[]string{"/deis/database/host", "/deis/database/port", "/deis/database/engine"},
		[]string{ipaddr, dbPort, "postgresql_psycopg2"}, dbhandler.C)
	if err != nil {
		t.Fatal(err)
	}
}
