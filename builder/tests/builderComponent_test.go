package verbose

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisBuilderTest(t *testing.T, testSessionUID string, port string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/builder:"+testSessionUID)
	dockercliutils.RunDeisDataTest(t, "--name", "deis-builder-data",
		"-v", "/var/lib/docker", "deis/base", "/bin/true")
	//docker run --name deis-builder -p 2223:22 -e PUBLISH=22 -e HOST=${COREOS_PRIVATE_IPV4} -e PORT=2223 --volumes-from deis-builder-data --privileged deis/builder
	IPAddress := utils.GetHostIPAddress()
	done <- true
	go func() {
		<-done
		dockercliutils.RunContainer(t, cli, "--name",
			"deis-builder-"+testSessionUID, "-p", "2223:22", "-e", "PUBLISH=22",
			"-e", "STORAGE_DRIVER=aufs", "-e", "HOST="+IPAddress, "-e",
			"ETCD_PORT="+port, "-e", "PORT=2223", "--volumes-from",
			"deis-builder-data", "--privileged", "deis/builder:"+testSessionUID)
	}()
	time.Sleep(5000 * time.Millisecond)
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-builder running")
}

func TestBuild(t *testing.T) {
	setkeys := []string{"/deis/registry/protocol",
		"deis/registry/host",
		"/deis/registry/port",
		"/deis/cache/host",
		"/deis/cache/port"}
	setdir := []string{"/deis/controller",
		"/deis/cache",
		"/deis/database",
		"/deis/registry",
		"/deis/domains"}
	var testSessionUID = utils.NewUuid()
	fmt.Println("UUID for the session Builder Test :" + testSessionUID)
	port := utils.GetRandomPort()
	//testSessionUID := "352aea64"
	dockercliutils.RunEtcdTest(t, testSessionUID, port)
	Builderhandler := etcdutils.InitetcdValues(setdir, setkeys, port)
	etcdutils.Publishvalues(t, Builderhandler)
	fmt.Println("starting Builder Component test")
	runDeisBuilderTest(t, testSessionUID, port)
	dockercliutils.DeisServiceTest(
		t, "deis-builder-"+testSessionUID, "22", "tcp")
	dockercliutils.ClearTestSession(t, testSessionUID)
}
