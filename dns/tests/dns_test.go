package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/utils"
)

func TestLogger(t *testing.T) {
	var err error
	tag, etcdPort := utils.BuildTag(), utils.RandomPort()

	//start etcd container
	etcdName := "deis-etcd-" + tag
	cli, stdout, stdoutPipe := dockercli.NewClient()
	dockercli.RunTestEtcd(t, etcdName, etcdPort)
	defer cli.CmdRm("-f", etcdName)

	host, port := utils.HostAddress(), utils.RandomPort()
	fmt.Printf("--- Run deis/dns:%s at %s:%s\n", tag, host, port)
	name := "deis-dns-" + tag
	defer cli.CmdRm("-f", name)
	go func() {
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-p", port+":53/udp",
			"-p", port+":53/tcp",
			"-e", "ETCD_MACHINES="+host+":"+etcdPort,
			"deis/dns:"+tag,
			"/bin/skydns",
			"-addr=0.0.0.0:53",
			"-nameservers=8.8.8.8:53,8.8.4.4:53",
			"-domain=deis.local.",
			"-verbose")
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "ready for queries on deis.local.")
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: Wait until etcd keys are published
	time.Sleep(5000 * time.Millisecond)
	dockercli.DeisServiceTest(t, name, port, "udp")
}
