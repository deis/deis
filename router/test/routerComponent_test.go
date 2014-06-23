package verbose

import (
  "fmt"
  "github.com/deis/deis/tests/dockercliutils"
  "github.com/deis/deis/tests/etcdutils"
  "github.com/deis/deis/tests/utils"
  "testing"
  "time"
)

func runDeisRouterTest(t *testing.T, testSessionUid string,port string) {
  cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
  done := make(chan bool, 1)
  dockercliutils.BuildDockerfile(t, "../","deis/router:"+testSessionUid)

  //ocker run --name deis-router -p 80:80 -p 2222:2222 -e PUBLISH=80 -e HOST=${COREOS_PRIVATE_IPV4} deis/router
  IPAddress :=  utils.GetHostIpAddress()
  done <- true
  go func() {
    <-done
    dockercliutils.RunContainer(t, cli,"--name", "deis-router-"+testSessionUid, "-p", "80:80","-p","2222:2222","-e","PUBLISH=80","-e", "HOST="+IPAddress,"-e","ETCD_PORT="+port, "deis/router:"+testSessionUid)
  }()
  time.Sleep(2000 * time.Millisecond)
  dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-router running")
}

func TestBuild(t *testing.T) {
  setkeys := []string{"deis/controller/host",
    "/deis/controller/port",
    "/deis/builder/host",
    "/deis/builder/port"}
  setdir := []string{"/deis/controller",
    "/deis/router",
    "/deis/database",
    "/deis/services",
    "/deis/builder",
    "/deis/domains"}
  var testSessionUid = utils.GetnewUuid()
  fmt.Println("UUID for the session Router Test :"+testSessionUid)
  port := utils.GetRandomPort()
  //testSessionUid := "352aea64"
  dockercliutils.RunEtcdTest(t, testSessionUid,port)
  Routerhandler := etcdutils.InitetcdValues(setdir, setkeys,port)
  etcdutils.Publishvalues(t, Routerhandler)
  fmt.Println("starting Router Component test")
  runDeisRouterTest(t, testSessionUid,port)
  dockercliutils.DeisServiceTest(t,"deis-router-"+testSessionUid,"80","http")
  dockercliutils.ClearTestSession(t, testSessionUid)
}
