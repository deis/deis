package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisLoggerTest(
	t *testing.T, testID string, etcdPort string, servicePort string) {
	var err error
	dockercliutils.RunDeisDataTest(t, "--name", "deis-logger-data",
		"-v", "/var/log/deis", "deis/base", "/bin/true")
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	go func() {
		err = dockercliutils.RunContainer(cli,
			"--name", "deis-logger-"+testID,
			"--rm",
			"-p", servicePort+":514/udp",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+utils.GetHostIPAddress(),
			"-e", "ETCD_PORT="+etcdPort,
			"--volumes-from", "deis-logger-data",
			"deis/logger:"+testID)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-logger running")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogger(t *testing.T) {
	testID := utils.NewUuid()
	err := dockercliutils.BuildImage(t, "../", "deis/logger:"+testID)
	if err != nil {
		t.Fatal(err)
	}
	etcdPort := utils.GetRandomPort()
	dockercliutils.RunEtcdTest(t, testID, etcdPort)
	servicePort := utils.GetRandomPort()
	fmt.Printf("--- Test deis-logger-%s at port %s\n", testID, servicePort)
	runDeisLoggerTest(t, testID, etcdPort, servicePort)
	dockercliutils.DeisServiceTest(
		t, "deis-logger-"+testID, servicePort, "udp")
	dockercliutils.ClearTestSession(t, testID)
}
