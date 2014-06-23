package verbose

import (
	"fmt"
	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/mockserviceutils"
	"github.com/deis/deis/tests/utils"
	"testing"
	"time"
)

func runDeisControllerTest(t *testing.T, testSessionUid string,port string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/controller:"+testSessionUid)
	//docker run --name deis-controller -p 8000:8000 -e PUBLISH=8000 -e HOST=${COREOS_PRIVATE_IPV4} --volumes-from=deis-logger deis/controller
	IPAddress :=  utils.GetHostIpAddress()
	done <- true
	go func() {
		<-done
		dockercliutils.RunContainer(t, cli, "--name", "deis-controller-"+testSessionUid, "-p", "8000:8000", "-e", "PUBLISH=8000", "-e", "HOST="+IPAddress,"-e","ETCD_PORT="+port, "deis/controller:"+testSessionUid)
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
	var testSessionUid = utils.GetnewUuid()
	fmt.Println("UUID for the session Controller Test :"+testSessionUid)
	port := utils.GetRandomPort()
	//testSessionUid := "352aea64"
	dockercliutils.RunEtcdTest(t, testSessionUid,port)
	fmt.Println("starting controller test:")
	Controllerhandler := etcdutils.InitetcdValues(setdir, setkeys,port)
	etcdutils.Publishvalues(t, Controllerhandler)
	mockserviceutils.RunMockDatabase(t, testSessionUid,port)
	fmt.Println("starting Controller component test")
	runDeisControllerTest(t, testSessionUid,port)
	dockercliutils.DeisServiceTest(t,"deis-controller-"+testSessionUid,"8000","http")
	dockercliutils.ClearTestSession(t, testSessionUid)
}
