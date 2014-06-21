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

func deisLoggerServiceTest(t *testing.T, testSessionUid string) {
	IPAddress := dockercliutils.GetInspectData(t, "{{ .NetworkSettings.IPAddress }}", "deis-logger-"+testSessionUid)
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

func runDeisLoggerTest(t *testing.T, testSessionUid string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/logger:"+testSessionUid)
	dockercliutils.RunDeisDataTest(t, "--name", "deis-logger-data", "-v", "/var/log/deis", "deis/base", "true")
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
		dockercliutils.RunContainer(t, cli, "--name", "deis-logger-"+testSessionUid, "-p", "514:514/udp", "-e", "PUBLISH=514", "-e", "HOST="+IPAddress, "--volumes-from", "deis-logger-data", "deis/logger:"+testSessionUid)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-logger running")
}

func TestBuild(t *testing.T) {
	fmt.Println("1st")
	var testSessionUid = utils.GetnewUuid()
	dockercliutils.RunDummyEtcdTest(t, testSessionUid)
	fmt.Println("starting registry test: %v", testSessionUid)
	fmt.Println("2nd")
	runDeisLoggerTest(t, testSessionUid)
	//deisRegistryServiceTest(t, testSessionUid)*/
	dockercliutils.ClearTestSession(t, testSessionUid)
	fmt.Println("3rd")

}
