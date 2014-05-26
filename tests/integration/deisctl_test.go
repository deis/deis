package integration

import (
    "fmt"
    "os/exec"
    "strings"
    "testing"

    "github.com/deis/deis/version"
)

func TestServerVersion(t *testing.T) {
    cmd := exec.Command("deisctl", "--version")
    output, err := cmd.CombinedOutput()

    if err != nil {
        t.Fatalf("Unexpected error while executing deisctl: %v", err)
    }

    if strings.TrimSpace(string(output)) != version.Version {
        t.Fatalf("Received unexpected output for `deisctl --version`: '%s'", output)
    }

    cmd = exec.Command("deisctl")
    output, err = cmd.CombinedOutput()

    if err != nil {
        t.Fatalf("Unexpected error while executing deisctl: %v", err)
    }

    if !strings.Contains(string(output), fmt.Sprintf("%s", version.Version)) {
        t.Fatalf("Could not find expected version string (%s) in help output:\n%s", version.Version, output)
    }
}
