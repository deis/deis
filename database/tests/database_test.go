package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/utils"
)

func runDeisDatabaseTest(t *testing.T, testID string, etcdPort string, servicePort string) {
	var err error
	dockercli.RunDeisDataTest(t, "--name", "deis-database-data",
		"-v", "/var/lib/postgresql", "deis/base", "true")
	cli, stdout, stdoutPipe := dockercli.GetNewClient()
	go func() {
		err = dockercli.RunContainer(cli,
			"--name", "deis-database-"+testID,
			"--rm",
			"-p", servicePort+":5432",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+utils.GetHostIPAddress(),
			"-e", "ETCD_PORT="+etcdPort,
			"--volumes-from", "deis-database-data",
			"deis/database:"+testID,
		)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "deis-database running")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDatabase(t *testing.T) {
	testID := utils.NewID()
	err := dockercli.BuildImage(t, "../", "deis/database:"+testID)
	if err != nil {
		t.Fatal(err)
	}
	etcdPort := utils.GetRandomPort()
	dockercli.RunEtcdTest(t, testID, etcdPort)
	servicePort := utils.GetRandomPort()
	fmt.Printf("--- Test deis-database-%s at port %s\n", testID, servicePort)
	runDeisDatabaseTest(t, testID, etcdPort, servicePort)
	dockercli.DeisServiceTest(t, "deis-database-"+testID, servicePort, "tcp")
	dockercli.ClearTestSession(t, testID)
}
