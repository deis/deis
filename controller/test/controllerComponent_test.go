package verbose

import (
	"fmt"
	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
	"net/http"
	"strings"
	"testing"
	"time"
)

func runDeisRegistryTest(t *testing.T, testSessionUid string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/registry:"+testSessionUid)
	dockercliutils.RunDeisDataTest(t, "--name", "deis-registry-data", "-v", "/data", "deis/base", "/bin/true")
	IPAddress := dockercliutils.GetInspectData(t, "{{ .NetworkSettings.IPAddress }}", "deis-etcd-"+testSessionUid)
	if strings.Contains(IPAddress, "Error") {
		t.Fatalf("worng IP %s", IPAddress)
	}
	done <- true
	go func() {
		<-done
		fmt.Println("inside run container")
		dockercliutils.RunContainer(t, cli, "--name", "deis-registry-"+testSessionUid, "-p", "5000:5000", "-e", "PUBLISH=5000", "-e", "HOST="+IPAddress, "--volumes-from", "deis-registry-data", "deis/registry:"+testSessionUid)
	}()
	time.Sleep(10000 * time.Millisecond)
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "Booting")

}

func deisRegistryServiceTest(t *testing.T, testSessionUid string) {
	IPAddress := dockercliutils.GetInspectData(t, "{{ .NetworkSettings.IPAddress }}", "deis-registry-"+testSessionUid)
	if strings.Contains(IPAddress, "Error") {
		t.Fatalf("worng IP %s", IPAddress)
	}
	url := "http://" + IPAddress + ":5000"
	response, err := http.Get(url)
	if err != nil {
		t.Fatalf("Not reachable %s", err)
	}
	fmt.Println(response)
}

func TestBuild(t *testing.T) {

	fmt.Println("1st")
	var testSessionUid = utils.GetnewUuid()
	testSessionUid = "352aea64"
	dockercliutils.RunEtcdTest(t, testSessionUid)
	fmt.Println("2nd")
	t.Logf("starting registry test: %v", testSessionUid)
	fmt.Println("starting registry test")
	runDeisRegistryTest(t, testSessionUid)
	deisRegistryServiceTest(t, testSessionUid)
	dockercliutils.ClearTestSession(t, testSessionUid)

}
