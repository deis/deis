package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisBuilderTest(t *testing.T, testID string, etcdPort string, servicePort string) {
	var err error
	dockercli.RunDeisDataTest(t, "--name", "deis-builder-data",
		"-v", "/var/lib/docker", "deis/base", "true")
	cli, stdout, stdoutPipe := dockercli.GetNewClient()
	go func() {
		err = dockercli.RunContainer(cli,
			"--name", "deis-builder-"+testID,
			"--rm",
			"-p", servicePort+":22",
			"-e", "PUBLISH=22",
			"-e", "STORAGE_DRIVER=aufs",
			"-e", "HOST="+utils.GetHostIPAddress(),
			"-e", "ETCD_PORT="+etcdPort,
			"-e", "PORT="+servicePort,
			"--volumes-from", "deis-builder-data",
			"--privileged", "deis/builder:"+testID,
		)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "deis-builder running")
	if err != nil {
		t.Fatal(err)
	}
}

func TestBuilder(t *testing.T) {
	setkeys := []string{
		"/deis/registry/protocol",
		"/deis/registry/host",
		"/deis/registry/port",
		"/deis/cache/host",
		"/deis/cache/port",
		"/deis/controller/protocol",
		"/deis/controller/host",
		"/deis/controller/port",
		"/deis/controller/builderKey",
	}
	setdir := []string{
		"/deis/controller",
		"/deis/cache",
		"/deis/database",
		"/deis/registry",
		"/deis/domains",
		"/deis/services",
	}
	testID := utils.NewID()
	err := dockercli.BuildImage(t, "../", "deis/builder:"+testID)
	if err != nil {
		t.Fatal(err)
	}
	etcdPort := utils.GetRandomPort()
	dockercli.RunEtcdTest(t, testID, etcdPort)
	handler := etcdutils.InitetcdValues(setdir, setkeys, etcdPort)
	etcdutils.Publishvalues(t, handler)
	servicePort := utils.GetRandomPort()
	fmt.Printf("--- Test deis-builder-%s at port %s\n", testID, servicePort)
	runDeisBuilderTest(t, testID, etcdPort, servicePort)
	// TODO: builder needs a few seconds to wake up here--fixme!
	time.Sleep(5000 * time.Millisecond)
	dockercli.DeisServiceTest(t, "deis-builder-"+testID, servicePort, "tcp")
	dockercli.ClearTestSession(t, testID)
}
