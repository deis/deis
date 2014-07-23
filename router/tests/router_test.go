package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisRouterTest(
	t *testing.T, testID string, etcdPort string, servicePort string) {
	var err error
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	go func() {
		err = dockercliutils.RunContainer(cli,
			"--name", "deis-router-"+testID,
			"--rm",
			"-p", servicePort+":80",
			"-p", "2222:2222",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+utils.GetHostIPAddress(),
			"-e", "ETCD_PORT="+etcdPort,
			"deis/router:"+testID)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-router running")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRouter(t *testing.T) {
	setkeys := []string{
		"deis/controller/host",
		"/deis/controller/port",
		"/deis/builder/host",
		"/deis/builder/port",
	}
	setdir := []string{
		"/deis/controller",
		"/deis/router",
		"/deis/database",
		"/deis/services",
		"/deis/builder",
		"/deis/domains",
	}
	testID := utils.NewUuid()
	err := dockercliutils.BuildImage(t, "../", "deis/router:"+testID)
	if err != nil {
		t.Fatal(err)
	}
	etcdPort := utils.GetRandomPort()
	dockercliutils.RunEtcdTest(t, testID, etcdPort)
	handler := etcdutils.InitetcdValues(setdir, setkeys, etcdPort)
	etcdutils.Publishvalues(t, handler)
	servicePort := utils.GetRandomPort()
	fmt.Printf("--- Test deis-router-%s at port %s\n", testID, servicePort)
	runDeisRouterTest(t, testID, etcdPort, servicePort)
	// TODO: nginx needs a few seconds to wake up here--fixme!
	time.Sleep(5000 * time.Millisecond)
	dockercliutils.DeisServiceTest(
		t, "deis-router-"+testID, servicePort, "http")
	dockercliutils.ClearTestSession(t, testID)
}
