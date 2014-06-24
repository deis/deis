package verbose

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisCacheTest(t *testing.T, testSessionUID string, etcdPort string, servicePort string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/cache:"+testSessionUID)
	IPAddress := utils.GetHostIPAddress()
	done <- true
	go func() {
		<-done
		//docker run --name deis-cache -p 6379:6379 -e PUBLISH=6379
		// -e HOST=${COREOS_PRIVATE_IPV4} deis/cache
		dockercliutils.RunContainer(t, cli, "--name",
			"deis-cache-"+testSessionUID,
			"-p", servicePort+":6379",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+IPAddress,
			"-e", "ETCD_PORT="+etcdPort,
			"deis/cache:"+testSessionUID)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "started")
}

func TestBuild(t *testing.T) {
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
