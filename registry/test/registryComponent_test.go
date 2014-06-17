package verbose

import (
	"github.com/deis/deis/tests/dockercliutils"
	"fmt"
	"testing"
)



func runDeisRegistryTest(t *testing.T){
	cli,stdout,stdoutPipe := dockercliutils.GetNewClient( )
	done := make(chan bool, 1)
	dockercliutils.BuildDockerfile(t,"../"," ")
	dockercliutils.RunDeisDataTest(t,"--name", "deis-registry-data", "-v", "/data", "deis/base", "/bin/true")
	IPAddress := dockercliutils.GetInspectData(t,"{{ .NetworkSettings.IPAddress }}", "deis-etcd")
	done <-true
	go func(){
		<- done
		fmt.Println("inside run etcd")
		dockercliutils.RunContainer(t,cli,"--name", "deis-registry", "-p", "5000:5000", "-e", "PUBLISH=5000", "-e", "HOST="+IPAddress, "--volumes-from", "deis-registry-data", "deis/registry")
	}()
	dockercliutils.PrintToStdout(t,stdout,stdoutPipe,"Booting")

}


func TestBuild(t *testing.T) {

	fmt.Println("1st")
	dockercliutils.RunEtcdTest(t)
	fmt.Println("2nd")

	fmt.Println("3rd")

	runDeisRegistryTest(t)


}
