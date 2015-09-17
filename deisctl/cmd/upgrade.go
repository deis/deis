package cmd

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/config"
	"github.com/deis/deis/deisctl/config/model"
)

// UpgradePrep stops and uninstalls all components except router and publisher
func UpgradePrep(stateless bool, b backend.Backend) error {
	var wg sync.WaitGroup

	b.Stop([]string{"database", "registry@*", "controller", "builder", "logger", "logspout"}, &wg, Stdout, Stderr)
	wg.Wait()
	b.Destroy([]string{"database", "registry@*", "controller", "builder", "logger", "logspout"}, &wg, Stdout, Stderr)
	wg.Wait()

	if !stateless {
		b.Stop([]string{"store-volume", "store-gateway@*"}, &wg, Stdout, Stderr)
		wg.Wait()
		b.Destroy([]string{"store-volume", "store-gateway@*"}, &wg, Stdout, Stderr)
		wg.Wait()

		b.Stop([]string{"store-metadata"}, &wg, Stdout, Stderr)
		wg.Wait()
		b.Destroy([]string{"store-metadata"}, &wg, Stdout, Stderr)
		wg.Wait()

		b.Stop([]string{"store-daemon"}, &wg, Stdout, Stderr)
		wg.Wait()
		b.Destroy([]string{"store-daemon"}, &wg, Stdout, Stderr)
		wg.Wait()

		b.Stop([]string{"store-monitor"}, &wg, Stdout, Stderr)
		wg.Wait()
		b.Destroy([]string{"store-monitor"}, &wg, Stdout, Stderr)
		wg.Wait()
	}

	fmt.Fprintln(Stdout, "The platform has been stopped, but applications are still serving traffic as normal.")
	fmt.Fprintln(Stdout, "Your cluster is now ready for upgrade. Install a new deisctl version and run `deisctl upgrade-takeover`.")
	fmt.Fprintln(Stdout, "For more details, see: http://docs.deis.io/en/latest/managing_deis/upgrading-deis/#graceful-upgrade")
	return nil
}

func listPublishedServices(cb config.Backend) ([]*model.ConfigNode, error) {
	nodes, err := cb.GetRecursive("deis/services")
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func republishServices(ttl uint64, nodes []*model.ConfigNode, cb config.Backend) error {
	for _, node := range nodes {
		_, err := cb.SetWithTTL(node.Key, node.Value, ttl)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpgradeTakeover gracefully starts a platform stopped with UpgradePrep
func UpgradeTakeover(stateless bool, b backend.Backend, cb config.Backend) error {
	if err := doUpgradeTakeOver(stateless, b, cb); err != nil {
		return err
	}

	return nil
}

func doUpgradeTakeOver(stateless bool, b backend.Backend, cb config.Backend) error {
	var wg sync.WaitGroup

	nodes, err := listPublishedServices(cb)
	if err != nil {
		return err
	}

	b.Stop([]string{"publisher"}, &wg, Stdout, Stderr)
	wg.Wait()
	b.Destroy([]string{"publisher"}, &wg, Stdout, Stderr)
	wg.Wait()

	if err := republishServices(1800, nodes, cb); err != nil {
		return err
	}

	b.RollingRestart("router", &wg, Stdout, Stderr)
	wg.Wait()
	b.Create([]string{"publisher"}, &wg, Stdout, Stderr)
	wg.Wait()
	b.Start([]string{"publisher"}, &wg, Stdout, Stderr)
	wg.Wait()

	installUpgradeServices(b, stateless, &wg, Stdout, Stderr)
	wg.Wait()

	startUpgradeServices(b, stateless, &wg, Stdout, Stderr)
	wg.Wait()
	return nil
}

func installUpgradeServices(b backend.Backend, stateless bool, wg *sync.WaitGroup, out, err io.Writer) {
	if !stateless {
		fmt.Fprintln(out, "Storage subsystem...")
		b.Create([]string{"store-daemon", "store-monitor", "store-metadata", "store-volume", "store-gateway@1"}, wg, out, err)
		wg.Wait()
	}

	fmt.Fprintln(out, "Logging subsystem...")
	if stateless {
		b.Create([]string{"logspout"}, wg, out, err)
	} else {
		b.Create([]string{"logger", "logspout"}, wg, out, err)
	}
	wg.Wait()

	fmt.Fprintln(out, "Control plane...")
	if stateless {
		b.Create([]string{"registry@1", "controller", "builder"}, wg, out, err)
	} else {
		b.Create([]string{"database", "registry@1", "controller", "builder"}, wg, out, err)
	}
	wg.Wait()

	fmt.Fprintln(out, "Data plane...")
	b.Create([]string{"publisher"}, wg, out, err)
	wg.Wait()
}

func startUpgradeServices(b backend.Backend, stateless bool, wg *sync.WaitGroup, out, err io.Writer) {

	// Wait for groups to come up.
	// If we're running in stateless mode, we start only a subset of services.
	if !stateless {
		fmt.Fprintln(out, "Storage subsystem...")
		b.Start([]string{"store-monitor"}, wg, out, err)
		wg.Wait()
		b.Start([]string{"store-daemon"}, wg, out, err)
		wg.Wait()
		b.Start([]string{"store-metadata"}, wg, out, err)
		wg.Wait()

		// we start gateway first to give metadata time to come up for volume
		b.Start([]string{"store-gateway@*"}, wg, out, err)
		wg.Wait()
		b.Start([]string{"store-volume"}, wg, out, err)
		wg.Wait()
	}

	// start logging subsystem first to collect logs from other components
	fmt.Fprintln(out, "Logging subsystem...")
	if !stateless {
		b.Start([]string{"logger"}, wg, out, err)
		wg.Wait()
	}
	b.Start([]string{"logspout"}, wg, out, err)
	wg.Wait()

	// Start these in parallel. This section can probably be removed now.
	var bgwg sync.WaitGroup
	var trash bytes.Buffer
	batch := []string{
		"database", "registry@*", "controller", "builder",
		"publisher",
	}
	if stateless {
		batch = []string{"registry@*", "controller", "builder", "publisher", "router@*"}
	}
	b.Start(batch, &bgwg, &trash, &trash)

	fmt.Fprintln(Stdout, "Control plane...")
	batch = []string{"database", "registry@*", "controller"}
	if stateless {
		batch = []string{"registry@*", "controller"}
	}
	b.Start(batch, wg, out, err)
	wg.Wait()

	b.Start([]string{"builder"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Data plane...")
	b.Start([]string{"publisher"}, wg, out, err)
	wg.Wait()
}
