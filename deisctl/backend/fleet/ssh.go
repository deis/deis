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
func (c *FleetClient) SSH(name string) (err error) {
	var sshClient *ssh.SSHForwardingClient

	timeout := time.Duration(Flags.SSHTimeout*1000) * time.Millisecond

	ms, err := machineState(name)
	if err != nil {
		return err
	}

	// If name isn't a machine ID, try it as a unit instead
	if ms == nil {
		units, err := c.Units(name)

		if err != nil {
			return err
		}

		machID, err := findUnit(units[0])

		if err != nil {
			return err
		}

		ms, err = machineState(machID)

		if err != nil || ms == nil {
			return err
		}
	}

	addr := ms.PublicIP

	if tun := getTunnelFlag(); tun != "" {
		sshClient, err = ssh.NewTunnelledSSHClient("core", tun, addr, getChecker(), false, timeout)
	} else {
		sshClient, err = ssh.NewSSHClient("core", addr, getChecker(), false, timeout)
	}
	if err != nil {
		return err
	}

	defer sshClient.Close()
	err = ssh.Shell(sshClient)
	return err
}

// runCommand will attempt to run a command on a given machine. It will attempt
// to SSH to the machine if it is identified as being remote.
func runCommand(cmd string, machID string) (retcode int) {
	var err error
	if machine.IsLocalMachineID(machID) {
		retcode, err = runLocalCommand(cmd)
		if err != nil {
			fmt.Printf("Error running local command: %v\n", err)
		}
	} else {
		ms, err := machineState(machID)
		if err != nil || ms == nil {
			fmt.Printf("Error getting machine IP: %v\n", err)
		} else {
			sshTimeout := time.Duration(Flags.SSHTimeout*1000) * time.Millisecond
			retcode, err = runRemoteCommand(cmd, ms.PublicIP, sshTimeout)
			if err != nil {
				fmt.Printf("Error running remote command: %v\n", err)
			}
		}
	}
	return
}

// runLocalCommand runs the given command locally and returns any error encountered and the exit code of the command
func runLocalCommand(cmd string) (int, error) {
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
func runRemoteCommand(cmd string, addr string, timeout time.Duration) (exit int, err error) {
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
func findUnit(name string) (machID string, err error) {
	u, err := cAPI.Unit(name)
	if err != nil {
		return "", fmt.Errorf("Error retrieving Unit %s: %v", name, err)
	}
	if suToGlobal(*u) {
		return "", fmt.Errorf("Unable to connect to global unit %s.\n", name)
	}
	if u == nil {
		return "", fmt.Errorf("Unit %s does not exist.\n", name)
	} else if u.CurrentState == "" {
		return "", fmt.Errorf("Unit %s does not appear to be running.\n", name)
	}

	return u.MachineID, nil
}

func machineState(machID string) (*machine.MachineState, error) {
	machines, err := cAPI.Machines()
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
func cachedMachineState(machID string) (ms *machine.MachineState) {
	if machineStates == nil {
		machineStates = make(map[string]*machine.MachineState)
		ms, err := cAPI.Machines()
		if err != nil {
			return nil
		}
		for i, m := range ms {
			machineStates[m.ID] = &ms[i]
		}
	}
	return machineStates[machID]
}
