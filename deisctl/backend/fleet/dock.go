package fleet

import (
	"fmt"
	"strings"
)

// Dock connects to the appropriate host and runs 'docker exec -it'.
func (c *FleetClient) Dock(target string, cmd []string) error {

	units, err := c.Units(target)
	if err != nil {
		return err
	}
	target = units[0][0 : len(units[0])-len(".service")]

	cmdStr := "sh"
	if len(cmd) > 0 {
		cmdStr = strings.Join(cmd, " ")
	}

	execit := fmt.Sprintf("docker exec -it %s %s", target, cmdStr)

	return c.SSHExec(target, execit)
}
