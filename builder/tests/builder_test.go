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
	t *testing.T, testSessionUID string, etcdPort string, servicePort string) {
	var err error
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.RunDeisDataTest(t, "--name", "deis-builder-data",
		"-v", "/var/lib/docker", "deis/base", "/bin/true")
	ipaddr := utils.GetHostIPAddress()
	done <- true
	go func() {
		<-done
		err = dockercliutils.RunContainer(cli,
			"--name", "deis-builder-"+testSessionUID,
			"--rm",
			"-p", servicePort+":22",
			"-e", "PUBLISH=22",
			"-e", "STORAGE_DRIVER=devicemapper",
			"-e", "HOST="+ipaddr,
			"-e", "ETCD_PORT="+etcdPort,
			"-e", "PORT="+servicePort,
			"--volumes-from", "deis-builder-data",
			"--privileged", "deis/builder:"+testSessionUID)
	}()
	time.Sleep(5000 * time.Millisecond)
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
	testSessionUID := utils.NewUuid()
	err := dockercliutils.BuildImage(t, "../", "deis/builder:"+testSessionUID)
	if err != nil {
		t.Fatal(err)
	}
	etcdPort := utils.GetRandomPort()
	dockercliutils.RunEtcdTest(t, testSessionUID, etcdPort)
	Builderhandler := etcdutils.InitetcdValues(setdir, setkeys, etcdPort)
	etcdutils.Publishvalues(t, Builderhandler)
	fmt.Println("starting Builder Component test")
	servicePort := utils.GetRandomPort()
	runDeisBuilderTest(t, testSessionUID, etcdPort, servicePort)
	// TODO: builder needs a few seconds to wake up here--fixme!
	time.Sleep(5000 * time.Millisecond)
	dockercliutils.DeisServiceTest(
		t, "deis-builder-"+testSessionUID, servicePort, "tcp")
	dockercliutils.ClearTestSession(t, testSessionUID)
}
