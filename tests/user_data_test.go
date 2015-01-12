// +build integration

package tests

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/deis/deis/tests/utils"
)

// check disable-transparent-huge-pages.service
func checkIfTHPAreDisabled(t *testing.T, cfg *utils.DeisTestConfig) {
	thp := "/sys/kernel/mm/transparent_hugepage/enabled"

	cmd := "sudo cat " + thp
	sshCmd := exec.Command("ssh",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "PasswordAuthentication=no",
		"core@deis."+cfg.Domain, cmd)
	out, err := sshCmd.Output()
	if err != nil {
		t.Fatal(out, err)
	}
	if !strings.Contains(string(out), "[never]") {
		t.Fatalf("Expected 'always madvise [never]' as selected value in %s but '%s' was returned", thp, string(out))
	}
}
