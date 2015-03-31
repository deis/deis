package confd

import (
	"fmt"
	"os"
	"time"

	Log "github.com/deis/deis/pkg/log"
	. "github.com/deis/deis/pkg/os"
)

var log = Log.New()

// WaitForInitialConf wait until the compilation of the templates is correct
func WaitForInitialConf(signalChan chan os.Signal, etcd string, timeout time.Duration) {
	log.Info("waiting for confd to write initial templates...")
	for {
		cmdAsString := fmt.Sprintf("confd -onetime -node %s -confdir /app", etcd)
		cmd, args := BuildCommandFromString(cmdAsString)
		err := RunCommand(signalChan, cmd, args, false)
		if err == nil {
			break
		}

		time.Sleep(timeout)
	}
}

// Launch launch confd as a daemon process.
func Launch(signalChan chan os.Signal, etcd string) {
	cmdAsString := fmt.Sprintf("confd -node %s -confdir /app --interval 5 --quiet --watch", etcd)
	cmd, args := BuildCommandFromString(cmdAsString)
	go RunProcessAsDaemon(signalChan, cmd, args)
}
