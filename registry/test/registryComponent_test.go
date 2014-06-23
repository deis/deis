package verbose

import (
	"fmt"
	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
	"testing"
	"time"
)

func runDeisRegistryTest(t *testing.T, testSessionUid string,port string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/registry:"+testSessionUid)
	dockercliutils.RunDeisDataTest(t, "--name", "deis-registry-data", "-v", "/data", "deis/base", "/bin/true")
	IPAddress :=  utils.GetHostIpAddress()
	done <- true
	go func() {
		<-done
		dockercliutils.RunContainer(t, cli, "--name", "deis-registry-"+testSessionUid, "-p", "5000:5000", "-e", "PUBLISH=5000", "-e", "HOST="+IPAddress,"-e","ETCD_PORT="+port,  "--volumes-from", "deis-registry-data", "deis/registry:"+testSessionUid)
	}()
	time.Sleep(2000 * time.Millisecond)
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "Booting")

}


func TestBuild(t *testing.T) {
	var testSessionUid = utils.GetnewUuid()
	fmt.Println("UUID for the session registry Test :"+testSessionUid)
	port := utils.GetRandomPort()
	//testSessionUid := "352aea64"
	dockercliutils.RunEtcdTest(t, testSessionUid,port)
	fmt.Println("starting registry component test")
	runDeisRegistryTest(t, testSessionUid,port)
	dockercliutils.DeisServiceTest(t,"deis-registry-"+testSessionUid,"5000","http")
	dockercliutils.ClearTestSession(t, testSessionUid)
}
