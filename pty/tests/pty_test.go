package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/utils"
)

func TestPty(t *testing.T) {
	var err error
	tag := utils.BuildTag()
	cli, stdout, stdoutPipe := dockercli.NewClient()
	host, port := utils.HostAddress(), utils.RandomPort()
	fmt.Printf("--- Run deis/pty:%s at %s:%s\n", tag, host, port)
	name := "deis-pty-" + tag
	defer cli.CmdRm("-f", name)
	go func() {
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-p", port+":3333",
			"-e", "COMMAND=ls",
			"deis/pty:"+tag)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "Listening on 0.0.0.0:3333")
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: Wait until etcd keys are published
	time.Sleep(5000 * time.Millisecond)
	dockercli.DeisServiceTest(t, name, port, "tcp")
}
