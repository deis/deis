package zookeeper

import (
	"io"
	"os/exec"
	"time"
)

// ZkServer struct to execute zookeeper commands.
type ZkServer struct {
	Stdout, Stderr io.Writer

	cmd *exec.Cmd
}

// Start starts a zookeeper server
func (srv *ZkServer) Start() error {
	srv.cmd = exec.Command("/opt/zookeeper/bin/zkServer.sh", "start-foreground")
	srv.cmd.Stdout = srv.Stdout
	srv.cmd.Stderr = srv.Stderr
	return srv.cmd.Start()
}

// Pid returns the process id of the running zookeeper server
func (srv *ZkServer) Pid() int {
	return srv.cmd.Process.Pid
}

// Stop stops a running zookeeper server
func (srv *ZkServer) Stop() {
	go func() {
		time.Sleep(1 * time.Second)
		srv.cmd.Process.Kill()
	}()
	srv.cmd.Process.Wait()
}
