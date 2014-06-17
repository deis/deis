package verbose

import (
  "github.com/deis/deis/tests/dockercliutils"
  "fmt"
  "testing"
)



func runDeisLoggerTest(t *testing.T){
  cli,stdout,stdoutPipe := dockercliutils.GetNewClient( )
  done := make(chan bool, 1)
  dockercliutils.BuildDockerfile(t,"../"," ")
  dockercliutils.RunDeisDataTest(t,"--name", "deis-logger-data", "-v", "/var/log/deis", "deis/base", "true")
  IPAddress := dockercliutils.GetInspectData(t,"{{ .NetworkSettings.IPAddress }}", "deis-etcd")
  done <-true
  go func(){
    <- done
    fmt.Println("inside run etcd")
    dockercliutils.RunContainer(t,cli,"--name", "deis-logger", "-p", "514:514/udp", "-e", "PUBLISH=514", "-e", "HOST="+IPAddress, "--volumes-from", "deis-logger-data", "deis/logger")
  }()
  dockercliutils.PrintToStdout(t,stdout,stdoutPipe,"Booting")

}


func TestBuild(t *testing.T) {

  fmt.Println("1st")
  dockercliutils.RunEtcdTest(t)
  fmt.Println("2nd")

  fmt.Println("3rd")

  runDeisLoggerTest(t)


}
