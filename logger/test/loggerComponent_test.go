package verbose

import (
	"fmt"
	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
	"testing"
)

func runDeisLoggerTest(t *testing.T, testSessionUid string,port string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/logger:"+testSessionUid)
	dockercliutils.RunDeisDataTest(t, "--name", "deis-logger-data", "-v", "/var/log/deis", "deis/base", "true")
	IPAddress :=  utils.GetHostIpAddress()
	done <- true
	go func() {
		<-done
		dockercliutils.RunContainer(t, cli, "--name", "deis-logger-"+testSessionUid, "-p", "514:514/udp", "-e", "PUBLISH=514", "-e", "HOST="+IPAddress, "-e","ETCD_PORT="+port, "--volumes-from", "deis-logger-data", "deis/logger:"+testSessionUid)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-logger running")
}

func TestBuild(t *testing.T) {
	var testSessionUid = utils.GetnewUuid()
	port := utils.GetRandomPort()
	fmt.Println("UUID for the session logger Test :"+testSessionUid)
	//testSessionUid := "352aea64"
	dockercliutils.RunEtcdTest(t, testSessionUid,port)
	fmt.Println("starting logger componenet test:")
	runDeisLoggerTest(t, testSessionUid,port)
	dockercliutils.DeisServiceTest(t,"deis-logger-"+testSessionUid,"514","udp")
	dockercliutils.ClearTestSession(t, testSessionUid)

}
