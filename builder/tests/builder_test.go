package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"
)

func TestBuilder(t *testing.T) {
	var err error
	setkeys := []string{
		"/deis/registry/protocol",
		"/deis/registry/host",
		"/deis/registry/port",
		"/deis/cache/host",
		"/deis/cache/port",
		"/deis/controller/protocol",
		"/deis/controller/host",
		"/deis/controller/port",
		"/deis/controller/builderKey",
	}
	setdir := []string{
		"/deis/controller",
		"/deis/cache",
		"/deis/database",
		"/deis/registry",
		"/deis/domains",
		"/deis/services",
	}
	tag, etcdPort := utils.BuildTag(), utils.RandomPort()
	etcdName := "deis-etcd-" + tag
	cli, stdout, stdoutPipe := dockercli.NewClient()
	dockercli.RunTestEtcd(t, etcdName, etcdPort)
	defer cli.CmdRm("-f", etcdName)
	handler := etcdutils.InitEtcd(setdir, setkeys, etcdPort)
	etcdutils.PublishEtcd(t, handler)
	ipaddr, port := utils.HostAddress(), utils.RandomPort()
	fmt.Printf("--- Run deis/builder:%s at %s:%s\n", tag, ipaddr, port)
	name := "deis-builder-" + tag
	defer cli.CmdRm("-f", name)
	go func() {
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-p", port+":22",
			"-e", "PUBLISH=22",
			"-e", "STORAGE_DRIVER=aufs",
			"-e", "HOST="+ipaddr,
			"-e", "ETCD_PORT="+etcdPort,
			"-e", "PORT="+port,
			"--privileged", "deis/builder:"+tag)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "deis-builder running")
	if err != nil {
		t.Fatal(err)
	}
	// TODO: builder needs a few seconds to wake up here--fixme!
	time.Sleep(5000 * time.Millisecond)
	dockercli.DeisServiceTest(t, name, port, "tcp")
}
