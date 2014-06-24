package verbose

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/mockserviceutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisControllerTest(t *testing.T, testSessionUID string, port string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/controller:"+testSessionUID)
	//docker run --name deis-controller -p 8000:8000 -e PUBLISH=8000 -e HOST=${COREOS_PRIVATE_IPV4} --volumes-from=deis-logger deis/controller
	IPAddress := utils.GetHostIPAddress()
	done <- true
	go func() {
		<-done
		dockercliutils.RunContainer(t, cli, "--name",
			"deis-controller-"+testSessionUID, "-p", "8000:8000",
			"-e", "PUBLISH=8000", "-e", "HOST="+IPAddress,
			"-e", "ETCD_PORT="+port, "deis/controller:"+testSessionUID)
	}()
	time.Sleep(5000 * time.Millisecond)
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "Booting")

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
	fmt.Println("UUID for the session Controller Test :" + testSessionUID)
	port := utils.GetRandomPort()
	//testSessionUID := "352aea64"
	dockercliutils.RunEtcdTest(t, testSessionUID, port)
	fmt.Println("starting controller test:")
	Controllerhandler := etcdutils.InitetcdValues(setdir, setkeys, port)
	etcdutils.Publishvalues(t, Controllerhandler)
	mockserviceutils.RunMockDatabase(t, testSessionUID, port)
	fmt.Println("starting Controller component test")
	runDeisControllerTest(t, testSessionUID, port)
	dockercliutils.DeisServiceTest(
		t, "deis-controller-"+testSessionUID, "8000", "http")
	dockercliutils.ClearTestSession(t, testSessionUID)
}
