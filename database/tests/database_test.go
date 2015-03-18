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

func TestDatabase(t *testing.T) {
	var err error
	tag, etcdPort := utils.BuildTag(), utils.RandomPort()
	imageName := utils.ImagePrefix() + "database" + ":" + tag
	cli, stdout, stdoutPipe := dockercli.NewClient()

	// start etcd container
	etcdName := "deis-etcd-" + tag
	dockercli.RunTestEtcd(t, etcdName, etcdPort)
	defer cli.CmdRm("-f", etcdName)

	// run mock ceph containers
	cephName := "deis-ceph-" + tag
	mock.RunMockCeph(t, cephName, cli, etcdPort)
	defer cli.CmdRm("-f", cephName)

	// run database container
	host, port := utils.HostAddress(), utils.RandomPort()
	fmt.Printf("--- Run %s at %s:%s\n", imageName, host, port)
	name := "deis-database-" + tag
	defer cli.CmdRm("-f", name)
	go func() {
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-p", port+":5432",
			"-e", "EXTERNAL_PORT="+port,
			"-e", "HOST="+host,
			"-e", "ETCD_PORT="+etcdPort,
			imageName)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "database: postgres is running...")
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: Wait until etcd keys are published
	time.Sleep(5000 * time.Millisecond)
	dockercli.DeisServiceTest(t, name, port, "tcp")
	etcdutils.VerifyEtcdValue(t, "/deis/database/host", host, etcdPort)
	etcdutils.VerifyEtcdValue(t, "/deis/database/port", port, etcdPort)
}
