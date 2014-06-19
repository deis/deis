package verbose

import (
	"fmt"
	"github.com/deis/deis/tests/dockercliutils"
	"strings"
	"testing"
)

func runDeisCacheTest(t *testing.T) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()

	done := make(chan bool, 1)

	dockercliutils.BuildDockerfile(t, "../", " ")
	IPAddress := dockercliutils.GetInspectData(t, "{{ .NetworkSettings.IPAddress }}", "deis-etcd")
	if strings.Contains(IPAddress, "Error") {
		t.Fatalf("worng IP %s", IPAddress)
	}

	fmt.Println(IPAddress + "IPADRESS")
	done <- true
	go func() {
		<-done
		fmt.Println("inside run cahce run continer")
		dockercliutils.RunContainer(t, cli, "--name", "deis-cache", "-p", "6379:6379", "-e", "PUBLISH=6379", "-e", "HOST="+IPAddress, "deis/cache")
	}()

	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "Server started")

}

func TestBuild(t *testing.T) {

	fmt.Println("1st")
	dockercliutils.RunEtcdTest(t, "foobar")
	fmt.Println("2nd")
	runDeisCacheTest(t)
}
