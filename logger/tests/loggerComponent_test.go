package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisLoggerTest(
	t *testing.T, testSessionUID string, etcdPort string, servicePort string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/logger:"+testSessionUID)
	dockercliutils.RunDeisDataTest(t, "--name", "deis-logger-data",
		"-v", "/var/log/deis", "deis/base", "true")
	IPAddress := utils.GetHostIPAddress()
	done <- true
	go func() {
		<-done
		dockercliutils.RunContainer(t, cli,
			"--name", "deis-logger-"+testSessionUID,
			"--rm",
			"-p", servicePort+":514/udp",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+IPAddress,
			"-e", "ETCD_PORT="+etcdPort,
			"--volumes-from", "deis-logger-data",
			"deis/logger:"+testSessionUID)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-logger running")
}

func TestBuild(t *testing.T) {
	var testSessionUID = utils.NewUuid()
	etcdPort := utils.GetRandomPort()
	servicePort := utils.GetRandomPort()
	fmt.Println("UUID for the session logger Test :" + testSessionUID)
	dockercliutils.RunEtcdTest(t, testSessionUID, etcdPort)
	fmt.Println("starting logger component test:")
	runDeisLoggerTest(t, testSessionUID, etcdPort, servicePort)
	dockercliutils.DeisServiceTest(
		t, "deis-logger-"+testSessionUID, servicePort, "udp")
	dockercliutils.ClearTestSession(t, testSessionUID)
}
