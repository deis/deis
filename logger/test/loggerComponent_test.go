package verbose

import (
	"fmt"
	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
	"net/http"
	"strings"
	"testing"
)

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

func runDeisLoggerTest(t *testing.T) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", " ")
	dockercliutils.RunDeisDataTest(t, "--name", "deis-logger-data", "-v", "/var/log/deis", "deis/base", "true")
	dockercliutils.RunDeisDataTest(t, "--name", "deis-logger-data", "-v", "/var/log/deis", "deis/base", "true")
	IPAddress := dockercliutils.GetInspectData(t, "{{ .NetworkSettings.IPAddress }}", "deis-etcd")
	if strings.Contains(IPAddress, "Error") {
		t.Fatalf("worng IP %s", IPAddress)
	}
	done <- true
	go func() {
		<-done
		fmt.Println("inside run etcd")
		dockercliutils.RunContainer(t, cli, "--name", "deis-logger", "-p", "514:514/udp", "-e", "PUBLISH=514", "-e", "HOST="+IPAddress, "--volumes-from", "deis-logger-data", "deis/logger")
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "Booting")
}

func TestBuild(t *testing.T) {

	fmt.Println("1st")
	var testSessionUid = utils.GetnewUuid()
	dockercliutils.RunEtcdTest(t)
	t.Logf("starting registry test: %v", testSessionUid)
	fmt.Println("2nd")
	runDeisLoggerTest(t, testSessionUid)
	deisRegistryServiceTest(t, testSessionUid)
	dockercliutils.ClearTestSession(t, testSessionUid)
	fmt.Println("3rd")

}
