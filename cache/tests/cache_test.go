package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/utils"
)

func runDeisCacheTest(t *testing.T, testID string, etcdPort string, servicePort string) {
	var err error
	cli, stdout, stdoutPipe := dockercli.GetNewClient()
	go func() {
		err = dockercli.RunContainer(cli,
			"--name", "deis-cache-"+testID,
			"--rm",
			"-p", servicePort+":6379",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+utils.GetHostIPAddress(),
			"-e", "ETCD_PORT="+etcdPort,
			"deis/cache:"+testID,
		)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "started")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCache(t *testing.T) {
	testID := utils.NewID()
	err := dockercli.BuildImage(t, "../", "deis/cache:"+testID)
	if err != nil {
		t.Fatal(err)
	}
	etcdPort := utils.GetRandomPort()
	dockercli.RunEtcdTest(t, testID, etcdPort)
	servicePort := utils.GetRandomPort()
	fmt.Printf("--- Test deis-cache-%s at port %s\n", testID, servicePort)
	runDeisCacheTest(t, testID, etcdPort, servicePort)
	dockercli.DeisServiceTest(t, "deis-cache-"+testID, servicePort, "tcp")
	dockercli.ClearTestSession(t, testID)
}
