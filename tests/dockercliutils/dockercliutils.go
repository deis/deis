package dockercliutils

import (
  "github.com/dotcloud/docker/api/client"
  "bufio"
  "fmt"
  "io"
  "time"
  "os"
  "testing"
  "strings"

)

const (
  unitTestStoreBase        = "/var/lib/docker/unit-tests"
  testDaemonAddr           = "172.17.8.100:4243"
  testDaemonProto          = "tcp"
  testDaemonHttpsProto     = "tcp"
)


func DaemonAddr() string {
  addr := os.Getenv("TEST_DAEMON_ADDR")
  if addr == "" {
    addr = "/var/run/docker.sock"
  }
  return addr
}

func DaemonProto() string {
  proto := os.Getenv("TEST_DAEMON_PROTO")
  if proto == "" {
    proto = "unix"
  }
  return proto
}

func CloseWrap(args ...io.Closer) error {
  e := false
  ret := fmt.Errorf("Error closing elements")
  for _, c := range args {
    if err := c.Close(); err != nil {
      e = true
      ret = fmt.Errorf("%s\n%s", ret, err)
    }
  }
  if e {
    return ret
  }
  return nil
}


func GetNewClient( ) (cli *client.DockerCli,stdout *io.PipeReader, stdoutPipe *io.PipeWriter){
  stdout, stdoutPipe = io.Pipe()
  cli = client.NewDockerCli(nil, stdoutPipe, nil, testDaemonProto, testDaemonAddr, nil)
  return //cli,stdout,stdoutpipe
}

func PrintToStdout( t *testing.T,stdout *io.PipeReader,stdoutPipe *io.PipeWriter,stoptag string) string {
  var result string
  for{
    if cmdBytes,err:= bufio.NewReader(stdout).ReadString('\n'); err==nil{
      result = cmdBytes
      fmt.Print(cmdBytes)
      if strings.Contains(cmdBytes, stoptag) == true {
        if err := CloseWrap(stdout, stdoutPipe); err != nil {
          t.Fatalf("Closewraps %s", err)
        }
      }
    }else{
      break
    }
  }
  return result
}


func BuildDockerfile(t *testing.T,path string,tag string){
  cli,stdout,stdoutPipe :=GetNewClient( )
  go func(){
    err:= cli.CmdBuild(path)
    if err != nil {
      t.Fatalf(" %s",err)
    }
    if err = CloseWrap(stdout, stdoutPipe); err != nil {
      t.Fatalf("buildCacheTest %s",err)
    }
  }( )
  time.Sleep(1000 * time.Millisecond)
  PrintToStdout(t ,stdout,stdoutPipe,"Building docker file")
}


func GetInspectData(t *testing.T,format string,container string) string{
  var inspectData string
  cli,stdout,stdoutPipe :=GetNewClient( )
  go func(){
    err:= cli.CmdInspect("--format",format,container)
    if err != nil {
      t.Fatalf("getIPAdressTest %s",err)
    }
    if err = CloseWrap(stdout, stdoutPipe); err != nil {
      t.Fatalf("getIPAdressTest %s",err)
    }
  }()
  go func() {
    fmt.Println("here")
    time.Sleep(3000 * time.Millisecond)
    if err := CloseWrap(stdout, stdoutPipe); err != nil {
      t.Fatalf("Inspect Element %s",err)
    }
  } ( )
  time.Sleep(1000 * time.Millisecond)
  inspectData = PrintToStdout(t ,stdout,stdoutPipe,"IPAddress")
  return strings.TrimSuffix(inspectData,"\n")

}



func PullImage(t *testing.T, cli *client.DockerCli,args ...string){
  err:= cli.CmdPull(args...)
  if err != nil {
    t.Fatalf("pulling Image Failed %s",err)
  }
}


func RunContainer(t *testing.T, cli *client.DockerCli,args ...string){
  err:= cli.CmdRun(args...)
  if err != nil {
    t.Fatalf("running Image failed %s",err)
  }
}


func RunDeisDataTest(t *testing.T,args ...string) {

  cli,stdout,stdoutPipe :=GetNewClient( )
  var hostname string
  go func() {
    hostname=GetInspectData(t,"{{ .Config.Hostname }}", args[1])
  }( )
  fmt.Println("inspecting deis registry data")

  if strings.Contains(hostname, "Error") == true {
    go func() {
      RunContainer(t,cli,args...)
    }()
     PrintToStdout(t ,stdout,stdoutPipe,"running"+args[1])
      fmt.Println("pulling Deis registry data")
    }
}



func RunEtcdTest(t *testing.T){
  cli,stdout,stdoutPipe :=GetNewClient( )
  done := make(chan bool, 1)
  done1 := make(chan bool,1)
  go func (){
    fmt.Println("inside pull etcd")
    PullImage(t,cli,"phife.atribecalledchris.com:5000/deis/etcd:0.3.0")
    done <- true
  }()
  go func (){
    <- done
    done1 <- true
    fmt.Println("inside run etcd")
    RunContainer(t,cli,"--name","deis-etcd","phife.atribecalledchris.com:5000/deis/etcd:0.3.0")
  }( )
  go func() {
    <-done1
    fmt.Println("here")
    time.Sleep(5000 * time.Millisecond)
    if err := CloseWrap(stdout, stdoutPipe); err != nil {
      t.Fatalf("runEtcdTest %s",err)
    }
  } ( )
  time.Sleep(1000 * time.Millisecond)
  PrintToStdout(t ,stdout,stdoutPipe,"pulling etcd")
}
