package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"
)

func TestStore(t *testing.T) {
	hostname := utils.Hostname()
	var err error

	// Set up etcd, which will be used by all containers
	tag, etcdPort := utils.BuildTag(), utils.RandomPort()

	etcdName := "deis-etcd-" + tag
	cli, stdout, stdoutPipe := dockercli.NewClient()
	dockercli.RunTestEtcd(t, etcdName, etcdPort)
	defer cli.CmdRm("-f", etcdName)
	host := utils.HostAddress()

	// prep etcd with the monitor hostname -- this is done in an ExecStartPre in the monitor unit
	etcdutils.SetSingle(t, "/deis/store/hosts/"+host, hostname, etcdPort)

	// since we're only running one OSD, our default of 128 placement groups is too large
	etcdutils.SetSingle(t, "/deis/store/pgNum", "64", etcdPort)

	// test deis-store-monitor
	imageName := utils.ImagePrefix() + "store-monitor" + ":" + tag
	fmt.Printf("--- Run %s at %s\n", imageName, host)
	name := "deis-store-monitor-" + tag
	defer cli.CmdRm("-f", name)
	go func() {
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-e", "HOST="+host,
			"-e", "ETCD_PORT="+etcdPort,
			"-e", "NUM_STORES=1",
			"--net=host",
			imageName)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "monmap e1: 1 mons at")
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: Wait until etcd keys are published
	time.Sleep(5000 * time.Millisecond)
	dockercli.DeisServiceTest(t, name, "6789", "tcp")
	etcdutils.VerifyEtcdKey(t, "/deis/store/monKeyring", etcdPort)
	etcdutils.VerifyEtcdKey(t, "/deis/store/adminKeyring", etcdPort)
	etcdutils.VerifyEtcdValue(t, "/deis/store/monSetupComplete", "youBetcha", etcdPort)

	// test deis-store-daemon
	imageName = utils.ImagePrefix() + "store-daemon" + ":" + tag
	fmt.Printf("--- Run %s at %s\n", imageName, host)
	name = "deis-store-daemon-" + tag
	cli2, stdout2, stdoutPipe2 := dockercli.NewClient()
	defer cli2.CmdRm("-f", "-v", name)
	go func() {
		_ = cli2.CmdRm("-f", "-v", name)
		err = dockercli.RunContainer(cli2,
			"--name", name,
			"--rm",
			"-e", "HOST="+host,
			"-e", "ETCD_PORT="+etcdPort,
			"--net=host",
			imageName)
	}()
	dockercli.PrintToStdout(t, stdout2, stdoutPipe2, "journal close /var/lib/ceph/osd/ceph-0/journal")
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: Wait until etcd keys are published
	time.Sleep(5000 * time.Millisecond)
	dockercli.DeisServiceTest(t, name, "6800", "tcp")
	etcdutils.VerifyEtcdValue(t, "/deis/store/osds/"+host, "0", etcdPort)

	// test deis-store-metadata
	imageName = utils.ImagePrefix() + "store-metadata" + ":" + tag
	fmt.Printf("--- Run %s at %s\n", imageName, host)
	name = "deis-store-metadata-" + tag
	cli3, stdout3, stdoutPipe3 := dockercli.NewClient()
	defer cli3.CmdRm("-f", "-v", name)
	go func() {
		_ = cli3.CmdRm("-f", "-v", name)
		err = dockercli.RunContainer(cli3,
			"--name", name,
			"--rm",
			"-e", "HOST="+host,
			"-e", "ETCD_PORT="+etcdPort,
			"--net=host",
			imageName)
	}()
	dockercli.PrintToStdout(t, stdout3, stdoutPipe3, "mds.0.1 active_start")
	if err != nil {
		t.Fatal(err)
	}

	// test deis-store-gateway
	imageName = utils.ImagePrefix() + "store-gateway" + ":" + tag
	port := utils.RandomPort()
	fmt.Printf("--- Run %s at %s:%s\n", imageName, host, port)
	name = "deis-store-gateway-" + tag
	cli4, stdout4, stdoutPipe4 := dockercli.NewClient()
	defer cli4.CmdRm("-f", name)
	go func() {
		_ = cli4.CmdRm("-f", name)
		err = dockercli.RunContainer(cli4,
			"--name", name,
			"--rm",
			"-h", "deis-store-gateway",
			"-p", port+":8888",
			"-e", "HOST="+host,
			"-e", "EXTERNAL_PORT="+port,
			"-e", "ETCD_PORT="+etcdPort,
			imageName)
	}()
	dockercli.PrintToStdout(t, stdout4, stdoutPipe4, "deis-store-gateway running...")
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: Wait until etcd keys are published
	time.Sleep(5000 * time.Millisecond)
	dockercli.DeisServiceTest(t, name, port, "http")
	etcdutils.VerifyEtcdValue(t, "/deis/store/gateway/host", host, etcdPort)
	etcdutils.VerifyEtcdValue(t, "/deis/store/gateway/port", port, etcdPort)
	etcdutils.VerifyEtcdKey(t, "/deis/store/gatewayKeyring", etcdPort)
	etcdutils.VerifyEtcdKey(t, "/deis/store/gateway/accessKey", etcdPort)
	etcdutils.VerifyEtcdKey(t, "/deis/store/gateway/secretKey", etcdPort)
}
