package verbose

import (
  "fmt"
  "github.com/deis/deis/tests/dockercliutils"
  "github.com/deis/deis/tests/etcdutils"
  "github.com/deis/deis/tests/mockserviceutils"
  "github.com/deis/deis/tests/utils"
  "net/http"
  "strings"
  "testing"
  "time"
)

func runDeisBuilderTest(t *testing.T, testSessionUid string) {
  cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
  done := make(chan bool, 1)
  dockercliutils.BuildDockerfile(t, "../", "deis/builder:"+testSessionUid)
  dockercliutils.RunDeisDataTest(t, "--name", "deis-builder-data", "-v", "/var/lib/docker", "deis/base", "/bin/true")
  //docker run --name deis-builder -p 2223:22 -e PUBLISH=22 -e HOST=${COREOS_PRIVATE_IPV4} -e PORT=2223 --volumes-from deis-builder-data --privileged deis/builder
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
    dockercliutils.RunContainer(t, cli, "--name", "deis-builder-"+testSessionUid, "-p", "2223:22", "-e", "PUBLISH=22", "-e", "STORAGE_DRIVER=aufs", "-e", "HOST="+IPAddress, "-e", "PORT=2223","--volumes-from","deis-builder-data","--privileged", "deis/builder:"+testSessionUid)
  }()
  time.Sleep(5000 * time.Millisecond)
  dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-builder running")
}

func deisBuilderServiceTest(t *testing.T, testSessionUid string) {
  IPAddress := dockercliutils.GetInspectData(t, "{{ .NetworkSettings.IPAddress }}", "deis-controller-"+testSessionUid)
  if strings.Contains(IPAddress, "Error") {
    t.Fatalf("worng IP %s", IPAddress)
  }
  url := "http://" + IPAddress + ":8000"
  response, err := http.Get(url)
  if err != nil {
    t.Fatalf("Not reachable %s", err)
  }
  fmt.Println(response)
}

func TestBuild(t *testing.T) {
  setkeys := []string{"/deis/registry/protocol",
    "deis/registry/host",
    "/deis/registry/port",
    "/deis/cache/host",
    "/deis/cache/port"}
  setdir := []string{"/deis/controller",
    "/deis/cache",
    "/deis/database",
    "/deis/registry",
    "/deis/domains"}
  fmt.Println("1st")
  var testSessionUid = utils.GetnewUuid()
  //testSessionUid := "352aea64"
  dockercliutils.RunDummyEtcdTest(t, testSessionUid)
  fmt.Println("2nd")
  t.Logf("starting controller test: %v", testSessionUid)
  Builderhandler := etcdutils.InitetcdValues(setdir, setkeys)
  etcdutils.Publishvalues(t, Builderhandler)
  fmt.Println("starting registry test")
  runDeisBuilderTest(t, testSessionUid)
  //deisBuilderServiceTest(t, testSessionUid)
  //dockercliutils.ClearTestSession(t, testSessionUid)
}
