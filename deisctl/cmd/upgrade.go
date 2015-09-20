package cmd

import (
	"fmt"
	"sync"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/config"
	"github.com/deis/deis/deisctl/config/model"
)

// UpgradePrep stops and uninstalls all components except router and publisher
func UpgradePrep(b backend.Backend) error {
	var wg sync.WaitGroup

	b.Stop([]string{"database", "registry@*", "controller", "builder", "logger", "logspout"}, &wg, Stdout, Stderr)
	wg.Wait()
	b.Destroy([]string{"database", "registry@*", "controller", "builder", "logger", "logspout"}, &wg, Stdout, Stderr)
	wg.Wait()

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
func UpgradeTakeover(b backend.Backend, cb config.Backend) error {

	if err := doUpgradeTakeOver(b, cb); err != nil {
		return err
	}

	return nil
}

func doUpgradeTakeOver(b backend.Backend, cb config.Backend) error {
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

	installDefaultServices(b, false, &wg, Stdout, Stderr) // @fixme: hax?
	wg.Wait()

	startDefaultServices(b, false, false, &wg, Stdout, Stderr) // @fixme: hax?
	wg.Wait()
	return nil
}
