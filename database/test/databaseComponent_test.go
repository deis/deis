package verbose

import (
	"fmt"
	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
	"testing"
)



func runDeisDatabaseTest(t *testing.T, testSessionUid string,port string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t, "../", "deis/database:"+testSessionUid)
	dockercliutils.RunDeisDataTest(t, "--name", "deis-database-data", "-v", "/var/lib/postgresql", "deis/base", "true")
	IPAddress :=  utils.GetHostIpAddress()
	done <- true
	go func() {
		<-done
		//docker run --name deis-database -p 5432:5432 -e PUBLISH=5432 -e HOST=${COREOS_PRIVATE_IPV4} --volumes-from deis-database-data deis/database
		dockercliutils.RunContainer(t, cli, "--name", "deis-database-"+testSessionUid, "-p", "5432:5432", "-e", "PUBLISH=5432", "-e", "HOST="+IPAddress,"-e","ETCD_PORT="+port, "--volumes-from", "deis-database-data", "deis/database:"+testSessionUid)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-database running")
}

func TestBuild(t *testing.T) {
	var testSessionUid = utils.GetnewUuid()
	fmt.Println("UUID for the session Cache Test :"+testSessionUid)
	port := utils.GetRandomPort()
	dockercliutils.RunEtcdTest(t, testSessionUid,port)
	fmt.Println("starting Database compotest:")
	runDeisDatabaseTest(t, testSessionUid,port)
	dockercliutils.DeisServiceTest(t,"deis-database-"+testSessionUid,"5432","tcp")
	dockercliutils.ClearTestSession(t, testSessionUid)
}
