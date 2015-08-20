// Package confd provides basic Confd support.
//
// Right now, this library is highly specific to the needs of the present
// builder. Because the confd library is not all public, we don't use it directly.
// Instead, we invoke the CLI.
package confd

import (
	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/cookoo/log"
	"github.com/Masterminds/cookoo/safely"

	"fmt"
	"os/exec"
	"strconv"
	"time"
)

// defaultEtcd is the default Etcd host.
const defaultEtcd = "127.0.0.1:4001"

// RunOnce runs the equivalent of `confd --onetime`.
//
// This may run the process repeatedly until either we time out (~20 minutes) or
// the templates are successfully built.
//
// Importantly, this blocks until the run is complete.
//
// Params:
// - node (string): The etcd node to use. (Only etcd is currently supported)
//
// Returns:
// - The []bytes from stdout and stderr when running the program.
//
func RunOnce(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {

	node := p.Get("node", defaultEtcd).(string)

	dargs := []string{"-onetime", "-node", node, "-log-level", "error"}

	log.Info(c, "Building confd templates. This may take a moment.")

	limit := 1200
	timeout := time.Second * 3
	var lasterr error
	start := time.Now()
	for i := 0; i < limit; i++ {
		if out, err := exec.Command("confd", dargs...).CombinedOutput(); err == nil {
			log.Infof(c, "Templates generated for %s on run %d", node, i)
			return out, nil
		} else {
			log.Debugf(c, "Recoverable error: %s", err)
			log.Debugf(c, "Output: %q", out)
			lasterr = err
		}

		time.Sleep(timeout)
		log.Infof(c, "Re-trying template build. (Elapsed time: %d)", time.Now().Sub(start)/time.Second)
	}

	return nil, fmt.Errorf("Could not build confd templates before timeout. Last error: %s", lasterr)
}

// Run starts confd and runs it in the background.
//
// If the command fails immediately on startup, an error is immediately
// returned. But from that point, a goroutine watches the command and
// reports if the command dies.
//
// Params:
// - node (string): The etcd node to use. (Only etcd is currently supported)
// - interval (int, default:5): The rebuilding interval.
//
// Returns
//  bool true if this succeeded.
func Run(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	// TODO: What should be done if confd dies at some point?

	node := p.Get("node", defaultEtcd).(string)
	interval := strconv.Itoa(p.Get("interval", 5).(int))

	cmd := exec.Command("confd", "-log-level", "error", "-node", node, "-interval", interval)
	if err := cmd.Start(); err != nil {
		return false, err
	}

	log.Infof(c, "Watching confd.")
	safely.Go(func() {
		if err := cmd.Wait(); err != nil {
			log.Infof(c, "confd exited with error: %s", err)
		}
	})

	return true, nil
}
