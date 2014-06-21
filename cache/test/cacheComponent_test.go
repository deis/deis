package verbose

import (
	"fmt"
	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
	"net/http"
	"strings"
	"testing"
)

func deisCacheServiceTest(t *testing.T, testSessionUid string) {
	IPAddress := dockercliutils.GetInspectData(t, "{{ .NetworkSettings.IPAddress }}", "deis-logger-"+testSessionUid)
	if strings.Contains(IPAddress, "Error") {
		t.Fatalf("worng IP %s", IPAddress)
	}
	url := "http://" + IPAddress + ":6379"
	response, err := http.Get(url)
	if err != nil {
		t.Fatalf("Not reachable %s", err)
	}
	fmt.Println(response)
}

func runDeisCacheTest(t *testing.T, testSessionUid string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/cache:"+testSessionUid)
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
		fmt.Println("inside run etcd")
		//docker run --name deis-cache -p 6379:6379 -e PUBLISH=6379 -e HOST=${COREOS_PRIVATE_IPV4} deis/cache
		dockercliutils.RunContainer(t, cli, "--name", "deis-cache-"+testSessionUid, "-p", "6379:6379", "-e", "PUBLISH=6379", "-e", "HOST="+IPAddress, "deis/cache:"+testSessionUid)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "started")
}

func TestBuild(t *testing.T) {
	fmt.Println("1st")
	var testSessionUid = utils.GetnewUuid()
	dockercliutils.RunDummyEtcdTest(t, testSessionUid)
	t.Logf("starting cache test: %v", testSessionUid)
	fmt.Println("2nd")
	runDeisCacheTest(t, testSessionUid)
	//deisCacheServiceTest(t, testSessionUid)
	//dockercliutils.ClearTestSession(t, testSessionUid)
	fmt.Println("3rd")

}
