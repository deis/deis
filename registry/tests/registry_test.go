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

func TestRegistry(t *testing.T) {
	var err error
	setkeys := []string{
		"/deis/cache/host",
		"/deis/cache/port",
		"/deis/store/gateway/host",
		"/deis/store/gateway/port",
		"/deis/store/gateway/accessKey",
		"/deis/store/gateway/secretKey",
	}
	setdir := []string{
		"/deis/cache",
		"/deis/store",
	}
	tag, etcdPort := utils.BuildTag(), utils.RandomPort()
	imageName := utils.ImagePrefix() + "registry" + ":" + tag
	etcdName := "deis-etcd-" + tag
	cli, stdout, stdoutPipe := dockercli.NewClient()
	dockercli.RunTestEtcd(t, etcdName, etcdPort)
	defer cli.CmdRm("-f", etcdName)
	handler := etcdutils.InitEtcd(setdir, setkeys, etcdPort)
	etcdutils.PublishEtcd(t, handler)

	// run mock ceph containers
	cephName := "deis-ceph-" + tag
	mock.RunMockCeph(t, cephName, cli, etcdPort)
	defer cli.CmdRm("-f", cephName)

	host, port := utils.HostAddress(), utils.RandomPort()
	fmt.Printf("--- Running %s at %s:%s\n", imageName, host, port)
	name := "deis-registry-" + tag
	defer cli.CmdRm("-f", name)
	go func() {
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-p", port+":5000",
			"-e", "EXTERNAL_PORT="+port,
			"-e", "HOST="+host,
			"-e", "ETCD_PORT="+etcdPort,
			imageName)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "listening on [::]:5000")
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: Wait until etcd keys are published
	time.Sleep(5000 * time.Millisecond)
	dockercli.DeisServiceTest(t, name, port, "http")
	etcdutils.VerifyEtcdValue(t, "/deis/registry/host", host, etcdPort)
	etcdutils.VerifyEtcdValue(t, "/deis/registry/port", port, etcdPort)
}
