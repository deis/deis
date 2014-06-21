package verbose

import (
	"fmt"
	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
	"net/http"
	"strings"
	"testing"
)

func deisDatabaseServiceTest(t *testing.T, testSessionUid string) {
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

func runDeisDatabaseTest(t *testing.T, testSessionUid string,port string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/database:"+testSessionUid)
	dockercliutils.RunDeisDataTest(t, "--name", "deis-database-data", "-v", "/var/lib/postgresql", "deis/base", "true")
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
		//docker run --name deis-database -p 5432:5432 -e PUBLISH=5432 -e HOST=${COREOS_PRIVATE_IPV4} --volumes-from deis-database-data deis/database
		dockercliutils.RunContainer(t, cli, "--name", "deis-database-"+testSessionUid, "-p", "5432:5432", "-e", "PUBLISH=5432", "-e", "HOST="+IPAddress,"-e","ETCD_PORT="+port, "--volumes-from", "deis-database-data", "deis/database:"+testSessionUid)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-database running")
}

func TestBuild(t *testing.T) {
	fmt.Println("1st")
	var testSessionUid = utils.GetnewUuid()
	port := utils.GetRandomPort()
	dockercliutils.RunDummyEtcdTest(t, testSessionUid,port)
	t.Logf("starting registry test: %v", testSessionUid)
	fmt.Println("2nd")
	runDeisDatabaseTest(t, testSessionUid,port)
	//deisRegistryServiceTest(t, testSessionUid)
	//dockercliutils.ClearTestSession(t, testSessionUid)
	fmt.Println("3rd")

}
