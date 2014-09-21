package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/utils"
)

func TestDatabase(t *testing.T) {
	var err error
	tag, etcdPort := utils.BuildTag(), utils.RandomPort()
	etcdName := "deis-etcd-" + tag
	cli, stdout, stdoutPipe := dockercli.NewClient()
	dockercli.RunTestEtcd(t, etcdName, etcdPort)
	defer cli.CmdRm("-f", etcdName)
	dockercli.RunDeisDataTest(t, "--name", "deis-database-data",
		"-v", "/var/lib/postgresql", "deis/base", "true")
	host, port := utils.HostAddress(), utils.RandomPort()
	fmt.Printf("--- Run deis/database:%s at %s:%s\n", tag, host, port)
	name := "deis-database-" + tag
	defer cli.CmdRm("-f", name)
	go func() {
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-p", port+":5432",
			"-e", "PUBLISH="+port,
			"-e", "HOST="+host,
			"-e", "ETCD_PORT="+etcdPort,
			"--volumes-from", "deis-database-data",
			"deis/database:"+tag)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "deis-database running")
	if err != nil {
		t.Fatal(err)
	}
	dockercli.DeisServiceTest(
		t, "deis-database-"+tag, port, "tcp")
}
