package commons

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"time"

	"github.com/deis/deis/boot/logger"
)

// Wait until the compilation of the templates
func WaitForInitialConfd(etcd string, timeout time.Duration) {
	for {
		var buffer bytes.Buffer
		output := bufio.NewWriter(&buffer)
		cmd := exec.Command("confd", "-onetime", "-node", etcd, "-config-file", "/app/confd.toml")
		cmd.Stdout = output
		cmd.Stderr = output
		err := cmd.Run()
		output.Flush()
		if err == nil {
			break
		}

		logger.Log.Info("waiting for confd to write initial templates...")
		logger.Log.Debugf("%s", buffer.String())
		time.Sleep(timeout)
	}
}

// Launch confd child
func LaunchConfd(etcd string) {
	cmd := exec.Command("confd", "-node", etcd, "-config-file", "/app/confd.toml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logger.Log.Warnf("confd terminated by error: %v", err)
	}
}
