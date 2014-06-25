package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisDatabaseTest(
	t *testing.T, testSessionUID string, etcdPort string, servicePort string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/database:"+testSessionUID)
	dockercliutils.RunDeisDataTest(t, "--name", "deis-database-data",
		"-v", "/var/lib/postgresql", "deis/base", "true")
	IPAddress := utils.GetHostIPAddress()
	done <- true
	go func() {
		<-done
		//docker run --name deis-database -p 5432:5432 -e PUBLISH=5432
		// -e HOST=${COREOS_PRIVATE_IPV4}
		// --volumes-from deis-database-data deis/database
		dockercliutils.RunContainer(t, cli, "--name",
			"deis-database-"+testSessionUID,
			"-p", servicePort+":5432",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+IPAddress,
			"-e", "ETCD_PORT="+etcdPort,
			"--volumes-from", "deis-database-data",
			"deis/database:"+testSessionUID)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-database running")
}

func TestBuild(t *testing.T) {
	var testSessionUID = utils.NewUuid()
	fmt.Println("UUID for the session Cache Test :" + testSessionUID)
	etcdPort := utils.GetRandomPort()
	servicePort := utils.GetRandomPort()
	dockercliutils.RunEtcdTest(t, testSessionUID, etcdPort)
	fmt.Println("starting Database component test:")
	runDeisDatabaseTest(t, testSessionUID, etcdPort, servicePort)
	dockercliutils.DeisServiceTest(
		t, "deis-database-"+testSessionUID, servicePort, "tcp")
	dockercliutils.ClearTestSession(t, testSessionUID)
}
