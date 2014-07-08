package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisCacheTest(t *testing.T, testSessionUID string, etcdPort string, servicePort string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	err := dockercliutils.BuildImage(t, "../", "deis/cache:"+testSessionUID)
	if err != nil {
		t.Fatal(err)
	}
	ipaddr := utils.GetHostIPAddress()
	done <- true
	go func() {
		<-done
		err = dockercliutils.RunContainer(cli,
			"--name", "deis-cache-"+testSessionUID,
			"--rm",
			"-p", servicePort+":6379",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+ipaddr,
			"-e", "ETCD_PORT="+etcdPort,
			"deis/cache:"+testSessionUID)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "started")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCache(t *testing.T) {
	var testSessionUID = utils.NewUuid()
	fmt.Println("UUID for the session Cache Test :" + testSessionUID)
	etcdPort := utils.GetRandomPort()
	servicePort := utils.GetRandomPort()
	dockercliutils.RunEtcdTest(t, testSessionUID, etcdPort)
	fmt.Println("starting cache component test:")
	runDeisCacheTest(t, testSessionUID, etcdPort, servicePort)
	dockercliutils.DeisServiceTest(
		t, "deis-cache-"+testSessionUID, servicePort, "tcp")
	dockercliutils.ClearTestSession(t, testSessionUID)
}
