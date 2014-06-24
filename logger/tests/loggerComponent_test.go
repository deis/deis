package verbose

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisLoggerTest(t *testing.T, testSessionUID string, port string) {
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
			"--name", "deis-logger-"+testSessionUID, "-p", "514:514/udp",
			"-e", "PUBLISH=514", "-e", "HOST="+IPAddress,
			"-e", "ETCD_PORT="+port, "--volumes-from", "deis-logger-data",
			"deis/logger:"+testSessionUID)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-logger running")
}

func TestBuild(t *testing.T) {
	var testSessionUID = utils.NewUuid()
	port := utils.GetRandomPort()
	fmt.Println("UUID for the session logger Test :" + testSessionUID)
	//testSessionUID := "352aea64"
	dockercliutils.RunEtcdTest(t, testSessionUID, port)
	fmt.Println("starting logger componenet test:")
	runDeisLoggerTest(t, testSessionUID, port)
	dockercliutils.DeisServiceTest(
		t, "deis-logger-"+testSessionUID, "514", "udp")
	dockercliutils.ClearTestSession(t, testSessionUID)
}
