package verbose

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisDatabaseTest(t *testing.T, testSessionUID string, port string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/database:"+testSessionUID)
	dockercliutils.RunDeisDataTest(t, "--name", "deis-database-data",
		"-v", "/var/lib/postgresql", "deis/base", "true")
	IPAddress := utils.GetHostIPAddress()
	done <- true
	go func() {
		<-done
		//docker run --name deis-database -p 5432:5432 -e PUBLISH=5432 -e HOST=${COREOS_PRIVATE_IPV4} --volumes-from deis-database-data deis/database
		dockercliutils.RunContainer(t, cli, "--name",
			"deis-database-"+testSessionUID, "-p", "5432:5432",
			"-e", "PUBLISH=5432", "-e", "HOST="+IPAddress,
			"-e", "ETCD_PORT="+port, "--volumes-from", "deis-database-data",
			"deis/database:"+testSessionUID)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-database running")
}

func TestBuild(t *testing.T) {
	var testSessionUID = utils.NewUuid()
	fmt.Println("UUID for the session Cache Test :" + testSessionUID)
	port := utils.GetRandomPort()
	dockercliutils.RunEtcdTest(t, testSessionUID, port)
	fmt.Println("starting Database compotest:")
	runDeisDatabaseTest(t, testSessionUID, port)
	dockercliutils.DeisServiceTest(
		t, "deis-database-"+testSessionUID, "5432", "tcp")
	dockercliutils.ClearTestSession(t, testSessionUID)
}
