package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisRegistryTest(
	t *testing.T, testID string, etcdPort string, servicePort string) {
	var err error
	dockercli.RunDeisDataTest(t, "--name", "deis-registry-data",
		"-v", "/data", "deis/base", "/bin/true")
	cli, stdout, stdoutPipe := dockercli.GetNewClient()
	go func() {
		err = dockercli.RunContainer(cli,
			"--name", "deis-registry-"+testID,
			"--rm",
			"-p", servicePort+":5000",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+utils.GetHostIPAddress(),
			"-e", "ETCD_PORT="+etcdPort,
			"--volumes-from", "deis-registry-data",
			"deis/registry:"+testID,
		)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "Booting")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegistry(t *testing.T) {
	setkeys := []string{
		"/deis/cache/host",
		"/deis/cache/port",
	}
	setdir := []string{
		"/deis/cache",
	}
	testID := utils.NewID()
	err := dockercli.BuildImage(t, "../", "deis/registry:"+testID)
	if err != nil {
		t.Fatal(err)
	}
	etcdPort := utils.GetRandomPort()
	dockercli.RunEtcdTest(t, testID, etcdPort)
	handler := etcdutils.InitetcdValues(setdir, setkeys, etcdPort)
	etcdutils.Publishvalues(t, handler)
	servicePort := utils.GetRandomPort()
	fmt.Printf("--- Test deis-registry-%s at port %s\n", testID, servicePort)
	runDeisRegistryTest(t, testID, etcdPort, servicePort)
	dockercli.DeisServiceTest(t, "deis-registry-"+testID, servicePort, "http")
	dockercli.ClearTestSession(t, testID)
}
