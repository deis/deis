package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisCacheTest(
	t *testing.T, testID string, etcdPort string, servicePort string) {
	var err error
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	go func() {
		err = dockercliutils.RunContainer(cli,
			"--name", "deis-cache-"+testID,
			"--rm",
			"-p", servicePort+":6379",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+utils.GetHostIPAddress(),
			"-e", "ETCD_PORT="+etcdPort,
			"deis/cache:"+testID)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "started")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCache(t *testing.T) {
	testID := utils.NewUuid()
	err := dockercliutils.BuildImage(t, "../", "deis/cache:"+testID)
	if err != nil {
		t.Fatal(err)
	}
	etcdPort := utils.GetRandomPort()
	dockercliutils.RunEtcdTest(t, testID, etcdPort)
	servicePort := utils.GetRandomPort()
	fmt.Printf("--- Test deis-cache-%s at port %s\n", testID, servicePort)
	runDeisCacheTest(t, testID, etcdPort, servicePort)
	dockercliutils.DeisServiceTest(
		t, "deis-cache-"+testID, servicePort, "tcp")
	dockercliutils.ClearTestSession(t, testID)
}
