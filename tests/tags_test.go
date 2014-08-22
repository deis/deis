// +build integration

package tests

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/deis/deis/tests/utils"
)

var (
	tagsListCmd  = "tags:list --app={{.AppName}}"
	tagsSetCmd   = "tags:set --app={{.AppName}} environ=test"
	tagsUnsetCmd = "tags:unset --app={{.AppName}} environ"
)

func tagsTest(t *testing.T, cfg *utils.DeisTestConfig, ver int) {
	configFleetMetadata(t, cfg)
	utils.Execute(t, tagsListCmd, cfg, false, "No tags defined")
	utils.Execute(t, tagsSetCmd, cfg, false, "test")
	utils.Execute(t, tagsListCmd, cfg, false, "test")
	utils.Execute(t, tagsUnsetCmd, cfg, false, "No tags defined")
}

// configFleetMetdata applies Fleet metadata configuration over SSH
// and restarts the Fleet systemd unit
func configFleetMetadata(t *testing.T, cfg *utils.DeisTestConfig) {
	// check for existing metadata configuration
	cmd := "sudo systemctl cat fleet.service"
	sshCmd := exec.Command("ssh",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "PasswordAuthentication=no",
		"core@deis."+cfg.Domain, cmd)
	out, err := sshCmd.Output()
	if err != nil {
		t.Fatal(out, err)
	}
	if strings.Contains(string(out), "FLEET_METADATA") {
		return
	}
	// append metadata to fleet unit
	metadata := "environ=test"
	cmd = fmt.Sprintf(`sudo /bin/sh -c 'echo Environment=\"FLEET_METADATA=%s\" >> /run/systemd/system/fleet.service.d/20-cloudinit.conf'`, metadata)
	sshCmd = exec.Command("ssh",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "PasswordAuthentication=no",
		"core@deis."+cfg.Domain, cmd)
	out, err = sshCmd.Output()
	if err != nil {
		t.Fatal(out, err)
	}
	// reload all units and restart fleet
	cmd = "sudo systemctl daemon-reload && sudo systemctl restart fleet"
	sshCmd = exec.Command("ssh",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "PasswordAuthentication=no",
		"core@deis."+cfg.Domain, cmd)
	out, err = sshCmd.Output()
	if err != nil {
		t.Fatal(out, err)
	}
	// take a nap while fleet restarts
	time.Sleep(5000 * time.Millisecond)
}
