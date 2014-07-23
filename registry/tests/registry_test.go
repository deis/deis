package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisRegistryTest(
	t *testing.T, testID string, etcdPort string, servicePort string) {
	var err error
	dockercliutils.RunDeisDataTest(t, "--name", "deis-registry-data",
		"-v", "/data", "deis/base", "/bin/true")
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	go func() {
		err = dockercliutils.RunContainer(cli,
			"--name", "deis-registry-"+testID,
			"--rm",
			"-p", servicePort+":5000",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+utils.GetHostIPAddress(),
			"-e", "ETCD_PORT="+etcdPort,
			"--volumes-from", "deis-registry-data",
			"deis/registry:"+testID)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "Booting")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegistry(t *testing.T) {
	testID := utils.NewUuid()
	err := dockercliutils.BuildImage(t, "../", "deis/registry:"+testID)
	if err != nil {
		t.Fatal(err)
	}
	etcdPort := utils.GetRandomPort()
	dockercliutils.RunEtcdTest(t, testID, etcdPort)
	servicePort := utils.GetRandomPort()
	fmt.Printf("--- Test deis-registry-%s at port %s\n", testID, servicePort)
	runDeisRegistryTest(t, testID, etcdPort, servicePort)
	dockercliutils.DeisServiceTest(
		t, "deis-registry-"+testID, servicePort, "http")
	dockercliutils.ClearTestSession(t, testID)
}
