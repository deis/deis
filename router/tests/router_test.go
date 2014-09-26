package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"
)

func TestRouter(t *testing.T) {
	var err error
	setkeys := []string{
		"deis/controller/host",
		"/deis/controller/port",
		"/deis/builder/host",
		"/deis/builder/port",
	}
	setdir := []string{
		"/deis/controller",
		"/deis/router",
		"/deis/database",
		"/deis/services",
		"/deis/builder",
		"/deis/domains",
	}
	tag, etcdPort := utils.BuildTag(), utils.RandomPort()
	etcdName := "deis-etcd-" + tag
	cli, stdout, stdoutPipe := dockercli.NewClient()
	dockercli.RunTestEtcd(t, etcdName, etcdPort)
	defer cli.CmdRm("-f", etcdName)
	handler := etcdutils.InitEtcd(setdir, setkeys, etcdPort)
	etcdutils.PublishEtcd(t, handler)
	host, port := utils.HostAddress(), utils.RandomPort()
	fmt.Printf("--- Run deis/router:%s at %s:%s\n", tag, host, port)
	name := "deis-router-" + tag
	go func() {
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-p", port+":80",
			"-p", utils.RandomPort()+":2222",
			"-e", "EXTERNAL_PORT="+port,
			"-e", "HOST="+host,
			"-e", "ETCD_PORT="+etcdPort,
			"deis/router:"+tag)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "deis-router running")
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: nginx needs a couple seconds to wake up here
	// FIXME: Wait until etcd keys are published
	time.Sleep(5000 * time.Millisecond)
	dockercli.DeisServiceTest(t, name, port, "http")
	routerKeyPrefix := "/deis/router/"+host
	etcdutils.VerifyEtcdValue(t, routerKeyPrefix+"/host", host, etcdPort)
	etcdutils.VerifyEtcdValue(t, routerKeyPrefix+"/port", port, etcdPort)
	_ = cli.CmdRm("-f", name)
}
