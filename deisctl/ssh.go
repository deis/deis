package deisctl

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/ssh"
)

// runCommand will attempt to run a command on a given machine
// It will attempt to SSH to the machine if it is identified as being remote.
func runCommand(cmd string, ms *machine.MachineState) (retcode int) {
	var err error
	if machine.IsLocalMachineState(ms) {
		err, retcode = runLocalCommand(cmd)
		if err != nil {
			fmt.Printf("Error running local command: %v\n", err)
		}
	} else {
		err, retcode = runRemoteCommand(cmd, ms.PublicIP)
		if err != nil {
			fmt.Printf("Error running remote command: %v\n", err)
		}
	}
	return
}

// runLocalCommand runs the given command locally
// and returns any error encountered and the exit code of the command
func runLocalCommand(cmd string) (error, int) {
	cmdSlice := strings.Split(cmd, " ")
	osCmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	osCmd.Stderr = os.Stderr
	osCmd.Stdout = os.Stdout
	osCmd.Start()
	err := osCmd.Wait()
	if err != nil {
		// Get the command's exit status if we can
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return nil, status.ExitStatus()
			}
		}
		// Otherwise, generic command error
		return err, -1
	}
	return nil, 0
}

// runRemoteCommand runs the given command over SSH on the given IP
// and returns any error encountered and the exit status of the command
func runRemoteCommand(cmd string, addr string) (err error, exit int) {
	var sshClient *ssh.SSHForwardingClient
	if tun := getTunnelFlag(); tun != "" {
		sshClient, err = ssh.NewTunnelledSSHClient("core", tun, addr, getChecker(), false)
	} else {
		sshClient, err = ssh.NewSSHClient("core", addr, getChecker(), false)
	}
	if err != nil {
		return err, -1
	}
	defer sshClient.Close()
	return ssh.Execute(sshClient, cmd)
}
