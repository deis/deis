package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisRouterTest(
	t *testing.T, testSessionID string, etcdPort string, servicePort string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	err := dockercliutils.BuildImage(t, "../", "deis/router:"+testSessionID)
	if err != nil {
		t.Fatal(err)
	}
	ipaddr := utils.GetHostIPAddress()
	done <- true
	go func() {
		<-done
		err = dockercliutils.RunContainer(cli,
			"--name", "deis-router-"+testSessionID,
			"--rm",
			"-p", servicePort+":80",
			"-p", "2222:2222",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+ipaddr,
			"-e", "ETCD_PORT="+etcdPort,
			"deis/router:"+testSessionID)
	}()
	time.Sleep(2000 * time.Millisecond)
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-router running")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRouter(t *testing.T) {
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
	var testSessionID = utils.NewUuid()
	fmt.Println("UUID for the session Router Test :" + testSessionID)
	etcdPort := utils.GetRandomPort()
	servicePort := utils.GetRandomPort()
	dockercliutils.RunEtcdTest(t, testSessionID, etcdPort)
	Routerhandler := etcdutils.InitetcdValues(setdir, setkeys, etcdPort)
	etcdutils.Publishvalues(t, Routerhandler)
	fmt.Println("starting Router Component test")
	runDeisRouterTest(t, testSessionID, etcdPort, servicePort)
	// TODO: nginx needs a few seconds to wake up here--fixme!
	time.Sleep(5000 * time.Millisecond)
	dockercliutils.DeisServiceTest(
		t, "deis-router-"+testSessionID, servicePort, "http")
	dockercliutils.ClearTestSession(t, testSessionID)
}
