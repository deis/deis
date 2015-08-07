package fleet

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/ssh"
)

// SSH opens an interactive shell to a machine in the cluster
func (c *FleetClient) SSH(name string) error {
	sshClient, _, err := c.sshConnect(name)
	if err != nil {
		return err
	}

	defer sshClient.Close()
	err = ssh.Shell(sshClient)
	return err
}

func (c *FleetClient) SSHExec(name, cmd string) error {

	conn, ms, err := c.sshConnect(name)
	if err != nil {
		return err
	}

	fmt.Printf("Executing '%s' on %s\n", cmd, ms.PublicIP)

	err, _ = ssh.Execute(conn, cmd)
	return err
}

func (c *FleetClient) sshConnect(name string) (*ssh.SSHForwardingClient, *machine.MachineState, error) {

	timeout := time.Duration(Flags.SSHTimeout*1000) * time.Millisecond

	ms, err := c.machineState(name)
	if err != nil {
		return nil, nil, err
	}

	// If name isn't a machine ID, try it as a unit instead
	if ms == nil {
		units, err := c.Units(name)

		if err != nil {
			return nil, nil, err
		}

		machID, err := c.findUnit(units[0])

		if err != nil {
			return nil, nil, err
		}

		ms, err = c.machineState(machID)

		if err != nil || ms == nil {
			return nil, nil, err
		}
	}

	addr := ms.PublicIP

	if tun := getTunnelFlag(); tun != "" {
		sshClient, err := ssh.NewTunnelledSSHClient("core", tun, addr, getChecker(), false, timeout)
		return sshClient, ms, err
	}
	sshClient, err := ssh.NewSSHClient("core", addr, getChecker(), false, timeout)
	return sshClient, ms, err
}

// runCommand will attempt to run a command on a given machine. It will attempt
// to SSH to the machine if it is identified as being remote.
func (c *FleetClient) runCommand(cmd string, machID string) (retcode int) {
	var err error
	if machine.IsLocalMachineID(machID) {
		retcode, err = c.runner.LocalCommand(cmd)
		if err != nil {
			fmt.Fprintf(c.errWriter, "Error running local command: %v\n", err)
		}
	} else {
		ms, err := c.machineState(machID)
		if err != nil || ms == nil {
			fmt.Fprintf(c.errWriter, "Error getting machine IP: %v\n", err)
		} else {
			sshTimeout := time.Duration(Flags.SSHTimeout*1000) * time.Millisecond
			retcode, err = c.runner.RemoteCommand(cmd, ms.PublicIP, sshTimeout)
			if err != nil {
				fmt.Fprintf(c.errWriter, "Error running remote command: %v\n", err)
			}
		}
	}
	return
}

type commandRunner interface {
	LocalCommand(string) (int, error)
	RemoteCommand(string, string, time.Duration) (int, error)
}

type sshCommandRunner struct{}

// runLocalCommand runs the given command locally and returns any error encountered and the exit code of the command
func (sshCommandRunner) LocalCommand(cmd string) (int, error) {
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
				return status.ExitStatus(), nil
			}
		}
		// Otherwise, generic command error
		return -1, err
	}
	return 0, nil
}

// runRemoteCommand runs the given command over SSH on the given IP, and returns
// any error encountered and the exit status of the command
func (sshCommandRunner) RemoteCommand(cmd string, addr string, timeout time.Duration) (exit int, err error) {
	var sshClient *ssh.SSHForwardingClient
	if tun := getTunnelFlag(); tun != "" {
		sshClient, err = ssh.NewTunnelledSSHClient("core", tun, addr, getChecker(), false, timeout)
	} else {
		sshClient, err = ssh.NewSSHClient("core", addr, getChecker(), false, timeout)
	}
	if err != nil {
		return -1, err
	}

	defer sshClient.Close()

	err, exit = ssh.Execute(sshClient, cmd)
	return
}

// findUnits returns the machine ID of a running unit
func (c *FleetClient) findUnit(name string) (machID string, err error) {
	u, err := c.Fleet.Unit(name)
	switch {
	case err != nil:
		return "", fmt.Errorf("Error retrieving Unit %s: %v", name, err)
	case suToGlobal(*u):
		return "", fmt.Errorf("Unable to connect to global unit %s.\n", name)
	case u == nil:
		return "", fmt.Errorf("Unit %s does not exist.\n", name)
	case u.CurrentState == "":
		return "", fmt.Errorf("Unit %s does not appear to be running.\n", name)
	}

	return u.MachineID, nil
}

func (c *FleetClient) machineState(machID string) (*machine.MachineState, error) {
	machines, err := c.Fleet.Machines()
	if err != nil {
		return nil, err
	}
	for _, ms := range machines {
		if ms.ID == machID {
			return &ms, nil
		}
	}
	return nil, nil
}

// cachedMachineState makes a best-effort to retrieve the MachineState of the given machine ID.
// It memoizes MachineState information for the life of a fleetctl invocation.
// Any error encountered retrieving the list of machines is ignored.
func (c *FleetClient) cachedMachineState(machID string) (ms *machine.MachineState) {
	if c.machineStates == nil {
		c.machineStates = make(map[string]*machine.MachineState)
		ms, err := c.Fleet.Machines()
		if err != nil {
			return nil
		}
		for i, m := range ms {
			c.machineStates[m.ID] = &ms[i]
		}
	}
	return c.machineStates[machID]
}
