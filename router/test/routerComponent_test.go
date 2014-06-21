package verbose

import (
  "fmt"
  "github.com/deis/deis/tests/dockercliutils"
  "github.com/deis/deis/tests/etcdutils"
  "github.com/deis/deis/tests/utils"
  "net/http"
  "strings"
  "testing"
  "time"
)

func runDeisRouterTest(t *testing.T, testSessionUid string,port string) {
  cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
  done := make(chan bool, 1)
  dockercliutils.BuildDockerfile(t, "../","deis/router:"+testSessionUid)

  //ocker run --name deis-router -p 80:80 -p 2222:2222 -e PUBLISH=80 -e HOST=${COREOS_PRIVATE_IPV4} deis/router
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
    fmt.Println("inside run container")
    dockercliutils.RunContainer(t, cli,"--name", "deis-router-"+testSessionUid, "-p", "80:80","-p","2222:2222","-e","PUBLISH=80","-e", "HOST="+IPAddress,"-e","ETCD_PORT="+port, "deis/router:"+testSessionUid)
  }()
  time.Sleep(5000 * time.Millisecond)
  dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-router running")
}

func deisRouterServiceTest(t *testing.T, testSessionUid string) {
  IPAddress := dockercliutils.GetInspectData(t, "{{ .NetworkSettings.IPAddress }}", "deis-controller-"+testSessionUid)
  if strings.Contains(IPAddress, "Error") {
    t.Fatalf("worng IP %s", IPAddress)
  }
  url := "http://" + IPAddress + ":80"
  response, err := http.Get(url)
  if err != nil {
    t.Fatalf("Not reachable %s", err)
  }
  fmt.Println(response)
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
  fmt.Println("1st")
  var testSessionUid = utils.GetnewUuid()
  port := utils.GetRandomPort()
  //testSessionUid := "352aea64"
  dockercliutils.RunDummyEtcdTest(t, testSessionUid,port)
  fmt.Println("2nd")
  t.Logf("starting controller test: %v", testSessionUid)
  Routerhandler := etcdutils.InitetcdValues(setdir, setkeys,port)
  etcdutils.Publishvalues(t, Routerhandler)
  fmt.Println("starting registry test")
  //mockserviceutils.RunMockDatabase(t, testSessionUid)
  runDeisRouterTest(t, testSessionUid,port)
  //deisBuilderServiceTest(t, testSessionUid)
  //dockercliutils.ClearTestSession(t, testSessionUid)
}
