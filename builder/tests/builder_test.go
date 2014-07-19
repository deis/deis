package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisBuilderTest(
	t *testing.T, testID string, etcdPort string, servicePort string) {
	var err error
	dockercliutils.RunDeisDataTest(t, "--name", "deis-builder-data",
		"-v", "/var/lib/docker", "deis/base", "true")
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	go func() {
		err = dockercliutils.RunContainer(cli,
			"--name", "deis-builder-"+testID,
			"--rm",
			"-p", servicePort+":22",
			"-e", "PUBLISH=22",
			"-e", "STORAGE_DRIVER=devicemapper",
			"-e", "HOST="+utils.GetHostIPAddress(),
			"-e", "ETCD_PORT="+etcdPort,
			"-e", "PORT="+servicePort,
			"--volumes-from", "deis-builder-data",
			"--privileged", "deis/builder:"+testID)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-builder running")
	if err != nil {
		t.Fatal(err)
	}
}

func TestBuilder(t *testing.T) {
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
	err := dockercliutils.BuildImage(t, "../", "deis/builder:"+testID)
	if err != nil {
		t.Fatal(err)
	}
	etcdPort := utils.GetRandomPort()
	dockercliutils.RunEtcdTest(t, testID, etcdPort)
	handler := etcdutils.InitetcdValues(setdir, setkeys, etcdPort)
	etcdutils.Publishvalues(t, handler)
	servicePort := utils.GetRandomPort()
	fmt.Printf("--- Test deis-builder-%s at port %s\n", testID, servicePort)
	runDeisBuilderTest(t, testID, etcdPort, servicePort)
	// TODO: builder needs a few seconds to wake up here--fixme!
	time.Sleep(5000 * time.Millisecond)
	dockercliutils.DeisServiceTest(
		t, "deis-builder-"+testID, servicePort, "tcp")
	dockercliutils.ClearTestSession(t, testID)
}
