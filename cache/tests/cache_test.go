package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"
)

func TestCache(t *testing.T) {
	var err error
	tag := utils.BuildTag()
	etcdPort := utils.RandomPort()
	etcdName := "deis-etcd-" + tag
	cli, stdout, stdoutPipe := dockercli.NewClient()
	dockercli.RunTestEtcd(t, etcdName, etcdPort)
	defer cli.CmdRm("-f", etcdName)
	host, port := utils.HostAddress(), utils.RandomPort()
	fmt.Printf("--- Run deis/cache:%s at %s:%s\n", tag, host, port)
	name := "deis-cache-" + tag
	defer cli.CmdRm("-f", name)
	go func() {
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-p", port+":6379",
			"-e", "EXTERNAL_PORT="+port,
			"-e", "HOST="+host,
			"-e", "ETCD_PORT="+etcdPort,
			"deis/cache:"+tag)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "started")
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: Wait until etcd keys are published
	time.Sleep(5000 * time.Millisecond)
	dockercli.DeisServiceTest(t, name, port, "tcp")
	etcdutils.VerifyEtcdValue(t, "/deis/cache/host", host, etcdPort)
	etcdutils.VerifyEtcdValue(t, "/deis/cache/port", port, etcdPort)
}
