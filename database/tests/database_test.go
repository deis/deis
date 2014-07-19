package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisDatabaseTest(
	t *testing.T, testID string, etcdPort string, servicePort string) {
	var err error
	dockercliutils.RunDeisDataTest(t, "--name", "deis-database-data",
		"-v", "/var/lib/postgresql", "deis/base", "true")
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	go func() {
		err = dockercliutils.RunContainer(cli,
			"--name", "deis-database-"+testID,
			"--rm",
			"-p", servicePort+":5432",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+utils.GetHostIPAddress(),
			"-e", "ETCD_PORT="+etcdPort,
			"--volumes-from", "deis-database-data",
			"deis/database:"+testID)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-database running")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDatabase(t *testing.T) {
	testID := utils.NewUuid()
	err := dockercliutils.BuildImage(t, "../", "deis/database:"+testID)
	if err != nil {
		t.Fatal(err)
	}
	etcdPort := utils.GetRandomPort()
	dockercliutils.RunEtcdTest(t, testID, etcdPort)
	servicePort := utils.GetRandomPort()
	fmt.Printf("--- Test deis-database-%s at port %s\n", testID, servicePort)
	runDeisDatabaseTest(t, testID, etcdPort, servicePort)
	dockercliutils.DeisServiceTest(
		t, "deis-database-"+testID, servicePort, "tcp")
	dockercliutils.ClearTestSession(t, testID)
}
