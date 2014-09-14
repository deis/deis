package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/utils"
)

func TestLogger(t *testing.T) {
	var err error
	tag, etcdPort := utils.BuildTag(), utils.RandomPort()
	etcdName := "deis-etcd-" + tag
	cli, stdout, stdoutPipe := dockercli.NewClient()
	dockercli.RunTestEtcd(t, etcdName, etcdPort)
	defer cli.CmdRm("-f", etcdName)
	dockercli.RunDeisDataTest(t, "--name", "deis-logger-data",
		"-v", "/var/log/deis", "deis/base", "/bin/true")
	ipaddr, port := utils.HostAddress(), utils.RandomPort()
	fmt.Printf("--- Run deis/logger:%s at %s:%s\n", tag, ipaddr, port)
	name := "deis-logger-" + tag
	defer cli.CmdRm("-f", name)
	go func() {
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-p", port+":514/udp",
			"-e", "PUBLISH="+port,
			"-e", "HOST="+utils.HostAddress(),
			"-e", "ETCD_PORT="+etcdPort,
			"--volumes-from", "deis-logger-data",
			"deis/logger:"+tag)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "deis-logger running")
	if err != nil {
		t.Fatal(err)
	}
	dockercli.DeisServiceTest(t, name, port, "udp")
}
