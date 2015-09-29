package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/mock"
	"github.com/deis/deis/tests/utils"
)

func TestController(t *testing.T) {
	var err error
	setkeys := []string{
		"/deis/registry/protocol",
		"/deis/registry/host",
		"/deis/registry/port",
		"/deis/platform/domain",
		"/deis/logs/host",
	}
	setdir := []string{
		"/deis/controller",
		"/deis/database",
		"/deis/registry",
		"/deis/domains",
		"/deis/scheduler",
	}
	tag, etcdPort := utils.BuildTag(), utils.RandomPort()
	imageName := utils.ImagePrefix() + "controller" + ":" + tag

	//start etcd container
	etcdName := "deis-etcd-" + tag
	cli, stdout, stdoutPipe := dockercli.NewClient()
	dockercli.RunTestEtcd(t, etcdName, etcdPort)
	defer cli.CmdRm("-f", etcdName)

	handler := etcdutils.InitEtcd(setdir, setkeys, etcdPort)
	etcdutils.PublishEtcd(t, handler)
	mock.RunMockDatabase(t, tag, etcdPort, utils.RandomPort())
	defer cli.CmdRm("-f", "deis-test-database-"+tag)
	host, port := utils.HostAddress(), utils.RandomPort()
	fmt.Printf("--- Run %s at %s:%s\n", imageName, host, port)
	name := "deis-controller-" + tag
	defer cli.CmdRm("-f", name)
	go func() {
		cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-v", "/var/run/docker.sock:/var/run/docker.sock",
			"-v", "/var/run/fleet.sock:/var/run/fleet.sock",
			"-p", port+":8000",
			"-e", "EXTERNAL_PORT="+port,
			"-e", "HOST="+host,
			"-e", "ETCD_PORT="+etcdPort,
			imageName)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "Booting")
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: Wait until etcd keys are published
	time.Sleep(5000 * time.Millisecond)
	dockercli.DeisServiceTest(t, name, port, "http")
	etcdutils.VerifyEtcdValue(t, "/deis/controller/host", host, etcdPort)
	etcdutils.VerifyEtcdValue(t, "/deis/controller/port", port, etcdPort)
}
