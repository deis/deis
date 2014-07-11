package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisRegistryTest(
	t *testing.T, testSessionUID string, etcdPort string, servicePort string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	err := dockercliutils.BuildImage(t, "../", "deis/registry:"+testSessionUID)
	if err != nil {
		t.Fatal(err)
	}
	dockercliutils.RunDeisDataTest(t, "--name", "deis-registry-data",
		"-v", "/data", "deis/base", "/bin/true")
	ipaddr := utils.GetHostIPAddress()
	done <- true
	go func() {
		<-done
		err = dockercliutils.RunContainer(cli,
			"--name", "deis-registry-"+testSessionUID,
			"--rm",
			"-p", servicePort+":5000",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+ipaddr,
			"-e", "ETCD_PORT="+etcdPort,
			"--volumes-from", "deis-registry-data",
			"deis/registry:"+testSessionUID)
	}()
	time.Sleep(2000 * time.Millisecond)
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "Booting")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegistry(t *testing.T) {
	var testSessionUID = utils.NewUuid()
	fmt.Println("UUID for the session registry Test :" + testSessionUID)
	etcdPort := utils.GetRandomPort()
	servicePort := utils.GetRandomPort()
	dockercliutils.RunEtcdTest(t, testSessionUID, etcdPort)
	fmt.Println("starting registry component test")
	runDeisRegistryTest(t, testSessionUID, etcdPort, servicePort)
	dockercliutils.DeisServiceTest(
		t, "deis-registry-"+testSessionUID, servicePort, "http")
	dockercliutils.ClearTestSession(t, testSessionUID)
}
