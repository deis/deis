package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/mockserviceutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisControllerTest(
	t *testing.T, testID string, etcdPort string, servicePort string) {
	var err error
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	go func() {
		err = dockercliutils.RunContainer(cli,
			"--name", "deis-controller-"+testID,
			"--rm",
			"-p", servicePort+":8000",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+utils.GetHostIPAddress(),
			"-e", "ETCD_PORT="+etcdPort,
			"deis/controller:"+testID)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "Booting")
	if err != nil {
		t.Fatal(err)
	}
}

func TestController(t *testing.T) {
	setkeys := []string{
		"/deis/registry/protocol",
		"deis/registry/host",
		"/deis/registry/port",
		"/deis/cache/host",
		"/deis/cache/port",
	}
	setdir := []string{
		"/deis/controller",
		"/deis/cache",
		"/deis/database",
		"/deis/registry",
		"/deis/domains",
	}
	testID := utils.NewUuid()
	err := dockercliutils.BuildImage(t, "../", "deis/controller:"+testID)
	if err != nil {
		t.Fatal(err)
	}
	etcdPort := utils.GetRandomPort()
	dockercliutils.RunEtcdTest(t, testID, etcdPort)
	handler := etcdutils.InitetcdValues(setdir, setkeys, etcdPort)
	etcdutils.Publishvalues(t, handler)
	dbPort := utils.GetRandomPort()
	mockserviceutils.RunMockDatabase(t, testID, etcdPort, dbPort)
	servicePort := utils.GetRandomPort()
	fmt.Printf("--- Test deis-controller-%s at port %s\n", testID, servicePort)
	runDeisControllerTest(t, testID, etcdPort, servicePort)
	dockercliutils.DeisServiceTest(
		t, "deis-controller-"+testID, servicePort, "http")
	dockercliutils.ClearTestSession(t, testID)
}
