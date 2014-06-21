package verbose

import (
	"fmt"
	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/mockserviceutils"
	"github.com/deis/deis/tests/utils"
	"net/http"
	"strings"
	"testing"
	"time"
)

func runDeisControllerTest(t *testing.T, testSessionUid string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/controller:"+testSessionUid)
	//docker run --name deis-controller -p 8000:8000 -e PUBLISH=8000 -e HOST=${COREOS_PRIVATE_IPV4} --volumes-from=deis-logger deis/controller
	IPAddress := func() string {
		var Ip string
		if utils.GetHostOs() == "darwin" {
			Ip = "172.17.8.100"
		}
		return Ip
	}()
	done <- true
	go func() {
		<-done
		fmt.Println("inside run container")
		dockercliutils.RunContainer(t, cli, "--name", "deis-controller-"+testSessionUid, "-p", "8000:8000", "-e", "PUBLISH=8000", "-e", "HOST="+IPAddress, "deis/controller:"+testSessionUid)
	}()
	time.Sleep(5000 * time.Millisecond)
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "Booting")

}

func deisControllerServiceTest(t *testing.T, testSessionUid string) {
	IPAddress := dockercliutils.GetInspectData(t, "{{ .NetworkSettings.IPAddress }}", "deis-controller-"+testSessionUid)
	if strings.Contains(IPAddress, "Error") {
		t.Fatalf("worng IP %s", IPAddress)
	}
	url := "http://" + IPAddress + ":8000"
	response, err := http.Get(url)
	if err != nil {
		t.Fatalf("Not reachable %s", err)
	}
	fmt.Println(response)
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
	fmt.Println("1st")
	var testSessionUid = utils.GetnewUuid()
	//testSessionUid := "352aea64"
	dockercliutils.RunDummyEtcdTest(t, testSessionUid)
	fmt.Println("2nd")
	t.Logf("starting controller test: %v", testSessionUid)
	Controllerhandler := etcdutils.InitetcdValues(setdir, setkeys)
	etcdutils.Publishvalues(t, Controllerhandler)
	fmt.Println("starting registry test")
	mockserviceutils.RunMockDatabase(t, testSessionUid)
	runDeisControllerTest(t, testSessionUid)
	deisControllerServiceTest(t, testSessionUid)
	dockercliutils.ClearTestSession(t, testSessionUid)
}
