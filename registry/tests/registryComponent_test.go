package verbose

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisRegistryTest(t *testing.T, testSessionUID string, port string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/registry:"+testSessionUID)
	dockercliutils.RunDeisDataTest(t, "--name", "deis-registry-data",
		"-v", "/data", "deis/base", "/bin/true")
	IPAddress := utils.GetHostIPAddress()
	done <- true
	go func() {
		<-done
		dockercliutils.RunContainer(t, cli,
			"--name", "deis-registry-"+testSessionUID, "-p", "5000:5000",
			"-e", "PUBLISH=5000", "-e", "HOST="+IPAddress,
			"-e", "ETCD_PORT="+port, "--volumes-from", "deis-registry-data",
			"deis/registry:"+testSessionUID)
	}()
	time.Sleep(2000 * time.Millisecond)
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "Booting")
}

func TestBuild(t *testing.T) {
	var testSessionUID = utils.NewUuid()
	fmt.Println("UUID for the session registry Test :" + testSessionUID)
	port := utils.GetRandomPort()
	//testSessionUID := "352aea64"
	dockercliutils.RunEtcdTest(t, testSessionUID, port)
	fmt.Println("starting registry component test")
	runDeisRegistryTest(t, testSessionUID, port)
	dockercliutils.DeisServiceTest(
		t, "deis-registry-"+testSessionUID, "5000", "http")
	dockercliutils.ClearTestSession(t, testSessionUID)
}
