package verbose

import (
	"fmt"
	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
	"testing"
)



func runDeisCacheTest(t *testing.T, testSessionUid string,port string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/cache:"+testSessionUid)
	IPAddress :=  utils.GetHostIpAddress()
	done <- true
	go func() {
		<-done
		//docker run --name deis-cache -p 6379:6379 -e PUBLISH=6379 -e HOST=${COREOS_PRIVATE_IPV4} deis/cache
		dockercliutils.RunContainer(t, cli, "--name", "deis-cache-"+testSessionUid, "-p", "6379:6379", "-e", "PUBLISH=6379", "-e", "HOST="+IPAddress,"-e","ETCD_PORT="+port,"deis/cache:"+testSessionUid)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "started")
}

func TestBuild(t *testing.T) {
	var testSessionUid = utils.GetnewUuid()
	fmt.Println("UUID for the session Cache Test :"+testSessionUid)
	port := utils.GetRandomPort()
	//testSessionUid := "352aea64"
	dockercliutils.RunEtcdTest(t, testSessionUid,port)
	fmt.Println("starting cache compotest:")
	runDeisCacheTest(t, testSessionUid,port)
	dockercliutils.DeisServiceTest(t,"deis-cache-"+testSessionUid,"6379","tcp")
	dockercliutils.ClearTestSession(t, testSessionUid)


}
