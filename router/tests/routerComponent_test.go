package verbose

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisRouterTest(t *testing.T, testSessionID string, port string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/router:"+testSessionID)

	//ocker run --name deis-router -p 80:80 -p 2222:2222 -e PUBLISH=80 -e HOST=${COREOS_PRIVATE_IPV4} deis/router
	IPAddress := utils.GetHostIPAddress()
	done <- true
	go func() {
		<-done
		dockercliutils.RunContainer(t, cli,
			"--name", "deis-router-"+testSessionID, "-p", "80:80",
			"-p", "2222:2222", "-e", "PUBLISH=80", "-e", "HOST="+IPAddress,
			"-e", "ETCD_PORT="+port, "deis/router:"+testSessionID)
	}()
	time.Sleep(2000 * time.Millisecond)
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-router running")
}

func TestBuild(t *testing.T) {
	setkeys := []string{"deis/controller/host",
		"/deis/controller/port",
		"/deis/builder/host",
		"/deis/builder/port"}
	setdir := []string{"/deis/controller",
		"/deis/router",
		"/deis/database",
		"/deis/services",
		"/deis/builder",
		"/deis/domains"}
	var testSessionID = utils.NewUuid()
	fmt.Println("UUID for the session Router Test :" + testSessionID)
	port := utils.GetRandomPort()
	//testSessionID := "352aea64"
	dockercliutils.RunEtcdTest(t, testSessionID, port)
	Routerhandler := etcdutils.InitetcdValues(setdir, setkeys, port)
	etcdutils.Publishvalues(t, Routerhandler)
	fmt.Println("starting Router Component test")
	runDeisRouterTest(t, testSessionID, port)
	dockercliutils.DeisServiceTest(
		t, "deis-router-"+testSessionID, "80", "http")
	dockercliutils.ClearTestSession(t, testSessionID)
}
