package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/config"
	"github.com/deis/deis/deisctl/utils"

	docopt "github.com/docopt/docopt-go"
)

const (
	// PlatformCommand is shorthand for "all the Deis components."
	PlatformCommand string = "platform"
	swarm           string = "swarm"
)

// ListUnits prints a list of installed units.
func ListUnits(argv []string, b backend.Backend) error {
	usage := `Prints a list of installed units.

Usage:
  deisctl list [options]
`
	// parse command-line arguments
	if _, err := docopt.Parse(usage, argv, true, "", false); err != nil {
		return err
	}
	return b.ListUnits()
}

// ListUnitFiles prints the contents of all defined unit files.
func ListUnitFiles(argv []string, b backend.Backend) error {
	err := b.ListUnitFiles()
	return err
}

// Scale grows or shrinks the number of running components.
// Currently "router", "registry" and "store-gateway" are the only types that can be scaled.
func Scale(argv []string, b backend.Backend) error {
	usage := `Grows or shrinks the number of running components.

Currently "router", "registry" and "store-gateway" are the only types that can be scaled.

Usage:
  deisctl scale [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}
	targets := args["<target>"].([]string)

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	for _, target := range targets {
		component, num, err := splitScaleTarget(target)
		if err != nil {
			return err
		}
		// the router is the only component that can scale at the moment
		if !strings.Contains(component, "router") && !strings.Contains(component, "registry") && !strings.Contains(component, "store-gateway") {
			return fmt.Errorf("cannot scale %s components", component)
		}
		b.Scale(component, num, &wg, outchan, errchan)
		wg.Wait()
	}
	close(outchan)
	return nil
}

// Start activates the specified components.
func Start(argv []string, b backend.Backend) error {
	usage := `Activates the specified components.

Usage:
  deisctl start [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	// if target is platform, install all services
	targets := args["<target>"].([]string)

	if len(targets) == 1 {
		if targets[0] == PlatformCommand {
			return StartPlatform(b)
		}
		if targets[0] == swarm {
			return StartSwarm(b)
		}
	}
	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	b.Start(targets, &wg, outchan, errchan)
	wg.Wait()
	close(outchan)

	return nil
}

// checkRequiredKeys exist in etcd
func checkRequiredKeys() error {
	if err := config.CheckConfig("/deis/platform/", "domain"); err != nil {
		return fmt.Errorf(`Missing platform domain, use:
deisctl config platform set domain=<your-domain>`)
	}

	if err := config.CheckConfig("/deis/platform/", "sshPrivateKey"); err != nil {
		fmt.Printf(`Warning: Missing sshPrivateKey, "deis run" will be unavailable. Use:
deisctl config platform set sshPrivateKey=<path-to-key>
`)
	}
	return nil
}

func startDefaultServices(b backend.Backend, wg *sync.WaitGroup, outchan chan string, errchan chan error) {

	// create separate channels for background tasks
	_outchan := make(chan string)
	_errchan := make(chan error)
	var _wg sync.WaitGroup

	// wait for groups to come up
	outchan <- fmt.Sprintf("Storage subsystem...")
	b.Start([]string{"store-monitor"}, wg, outchan, errchan)
	wg.Wait()
	b.Start([]string{"store-daemon"}, wg, outchan, errchan)
	wg.Wait()
	b.Start([]string{"store-metadata"}, wg, outchan, errchan)
	wg.Wait()

	// we start gateway first to give metadata time to come up for volume
	b.Start([]string{"store-gateway@*"}, wg, outchan, errchan)
	wg.Wait()
	b.Start([]string{"store-volume"}, wg, outchan, errchan)
	wg.Wait()

	// start logging subsystem first to collect logs from other components
	outchan <- fmt.Sprintf("Logging subsystem...")
	b.Start([]string{"logger"}, wg, outchan, errchan)
	wg.Wait()
	b.Start([]string{"logspout"}, wg, outchan, errchan)
	wg.Wait()

	b.Start([]string{
		"database", "registry@*", "controller", "builder",
		"publisher", "router@*"},
		&_wg, _outchan, _errchan)

	outchan <- fmt.Sprintf("Control plane...")
	b.Start([]string{"database", "registry@*", "controller"}, wg, outchan, errchan)
	wg.Wait()
	b.Start([]string{"builder"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Data plane...")
	b.Start([]string{"publisher"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Routing mesh...")
	b.Start([]string{"router@*"}, wg, outchan, errchan)
	wg.Wait()
}

// Stop deactivates the specified components.
func Stop(argv []string, b backend.Backend) error {
	usage := `Deactivates the specified components.

Usage:
  deisctl stop [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}
	targets := args["<target>"].([]string)

	// if target is platform, stop all services
	if len(targets) == 1 {
		if targets[0] == PlatformCommand {
			return StopPlatform(b)
		}
		if targets[0] == swarm {
			return StopSwarm(b)
		}
	}

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	b.Stop(targets, &wg, outchan, errchan)
	wg.Wait()
	close(outchan)

	return nil
}

func stopDefaultServices(b backend.Backend, wg *sync.WaitGroup, outchan chan string, errchan chan error) {

	outchan <- fmt.Sprintf("Routing mesh...")
	b.Stop([]string{"router@*"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Data plane...")
	b.Stop([]string{"publisher"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Control plane...")
	b.Stop([]string{"controller", "builder", "database", "registry@*"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Logging subsystem...")
	b.Stop([]string{"logger", "logspout"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Storage subsystem...")
	b.Stop([]string{"store-volume", "store-gateway@*"}, wg, outchan, errchan)
	wg.Wait()
	b.Stop([]string{"store-metadata"}, wg, outchan, errchan)
	wg.Wait()
	b.Stop([]string{"store-daemon"}, wg, outchan, errchan)
	wg.Wait()
	b.Stop([]string{"store-monitor"}, wg, outchan, errchan)
	wg.Wait()
}

// Restart stops and then starts the specified components.
func Restart(argv []string, b backend.Backend) error {
	usage := `Stops and then starts the specified components.

Usage:
  deisctl restart [<target>...] [options]
`
	// parse command-line arguments
	if _, err := docopt.Parse(usage, argv, true, "", false); err != nil {
		return err
	}

	// act as if the user called "stop" and then "start"
	argv[0] = "stop"
	if err := Stop(argv, b); err != nil {
		return err
	}
	argv[0] = "start"
	return Start(argv, b)
}

// Status prints the current status of components.
func Status(argv []string, b backend.Backend) error {
	usage := `Prints the current status of components.

Usage:
  deisctl status [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	targets := args["<target>"].([]string)
	for _, target := range targets {
		if err := b.Status(target); err != nil {
			return err
		}
	}
	return nil
}

// Journal prints log output for the specified components.
func Journal(argv []string, b backend.Backend) error {
	usage := `Prints log output for the specified components.

Usage:
  deisctl journal [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	targets := args["<target>"].([]string)
	for _, target := range targets {
		if err := b.Journal(target); err != nil {
			return err
		}
	}
	return nil
}

// Install loads the definitions of components from local unit files.
// After Install, the components will be available to Start.
func Install(argv []string, b backend.Backend) error {
	usage := `Loads the definitions of components from local unit files.

After install, the components will be available to start.

"deisctl install" looks for unit files in these directories, in this order:
- the $DEISCTL_UNITS environment variable, if set
- $HOME/.deis/units
- /var/lib/deis/units

Usage:
  deisctl install [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	targets := args["<target>"].([]string)
	// if target is platform, install all services
	if len(targets) == 1 {
		if targets[0] == PlatformCommand {
			return InstallPlatform(b)
		}
		if targets[0] == swarm {
			return InstallSwarm(b)
		}
	}
	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	// otherwise create the specific targets
	b.Create(targets, &wg, outchan, errchan)
	wg.Wait()

	close(outchan)
	return nil
}

func installDefaultServices(b backend.Backend, wg *sync.WaitGroup, outchan chan string, errchan chan error) {

	outchan <- fmt.Sprintf("Storage subsystem...")
	b.Create([]string{"store-daemon", "store-monitor", "store-metadata", "store-volume", "store-gateway@1"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Logging subsystem...")
	b.Create([]string{"logger", "logspout"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Control plane...")
	b.Create([]string{"database", "registry@1", "controller", "builder"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Data plane...")
	b.Create([]string{"publisher"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Routing mesh...")
	b.Create([]string{"router@1", "router@2", "router@3"}, wg, outchan, errchan)
	wg.Wait()
}

// Uninstall unloads the definitions of the specified components.
// After Uninstall, the components will be unavailable until Install is called.
func Uninstall(argv []string, b backend.Backend) error {
	usage := `Unloads the definitions of the specified components.

After uninstall, the components will be unavailable until install is called.

Usage:
  deisctl uninstall [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	// if target is platform, uninstall all services
	targets := args["<target>"].([]string)
	if len(targets) == 1 {
		if targets[0] == PlatformCommand {
			return UninstallPlatform(b)
		}
		if targets[0] == swarm {
			return UnInstallSwarm(b)
		}
	}

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	// uninstall the specific target
	b.Destroy(targets, &wg, outchan, errchan)
	wg.Wait()
	close(outchan)

	return nil
}

func uninstallAllServices(b backend.Backend, wg *sync.WaitGroup, outchan chan string, errchan chan error) error {

	outchan <- fmt.Sprintf("Routing mesh...")
	b.Destroy([]string{"router@*"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Data plane...")
	b.Destroy([]string{"publisher"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Control plane...")
	b.Destroy([]string{"controller", "builder", "database", "registry@*"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Logging subsystem...")
	b.Destroy([]string{"logger", "logspout"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Storage subsystem...")
	b.Destroy([]string{"store-volume", "store-gateway@*"}, wg, outchan, errchan)
	wg.Wait()
	b.Destroy([]string{"store-metadata"}, wg, outchan, errchan)
	wg.Wait()
	b.Destroy([]string{"store-daemon"}, wg, outchan, errchan)
	wg.Wait()
	b.Destroy([]string{"store-monitor"}, wg, outchan, errchan)
	wg.Wait()

	return nil
}

func printState(outchan chan string, errchan chan error, interval time.Duration) error {
	for {
		select {
		case out := <-outchan:
			// done on closed channel
			if out == "" {
				return nil
			}
			fmt.Println(out)
		case err := <-errchan:
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
		}
		time.Sleep(interval)
	}
}

func splitScaleTarget(target string) (c string, num int, err error) {
	r := regexp.MustCompile(`([a-z-]+)=([\d]+)`)
	match := r.FindStringSubmatch(target)
	if len(match) == 0 {
		err = fmt.Errorf("Could not parse: %v", target)
		return
	}
	c = match[1]
	num, err = strconv.Atoi(match[2])
	if err != nil {
		return
	}
	return
}

// Config gets or sets a configuration value from the cluster.
//
// A configuration value is stored and retrieved from a key/value store (in this case, etcd)
// at /deis/<component>/<config>. Configuration values are typically used for component-level
// configuration, such as enabling TLS for the routers.
func Config(argv []string) error {
	usage := `Gets or sets a configuration value from the cluster.

A configuration value is stored and retrieved from a key/value store
(in this case, etcd) at /deis/<component>/<config>. Configuration
values are typically used for component-level configuration, such as
enabling TLS for the routers.

Note: "deisctl config platform set sshPrivateKey=" expects a path
to a private key.

Usage:
  deisctl config <target> get [<key>...]
  deisctl config <target> set <key=val>...

Examples:
  deisctl config platform set domain=mydomain.com
  deisctl config platform set sshPrivateKey=$HOME/.ssh/deis
  deisctl config controller get webEnabled
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}
	if err := config.Config(args); err != nil {
		return err
	}
	return nil
}

// RefreshUnits overwrites local unit files with those requested.
// Downloading from the Deis project GitHub URL by tag or SHA is the only mechanism
// currently supported.
func RefreshUnits(argv []string) error {
	usage := `Overwrites local unit files with those requested.

Downloading from the Deis project GitHub URL by tag or SHA is the only mechanism
currently supported.

"deisctl install" looks for unit files in these directories, in this order:
- the $DEISCTL_UNITS environment variable, if set
- $HOME/.deis/units
- /var/lib/deis/units

Usage:
  deisctl refresh-units [-p <target>] [-t <tag>]

Options:
  -p --path=<target>   where to save unit files [default: $HOME/.deis/units]
  -t --tag=<tag>       git tag, branch, or SHA to use when downloading unit files
                       [default: master]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}
	dir := args["--path"].(string)
	dir = utils.ResolvePath(dir)
	// create the target dir if necessary
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	// download and save the unit files to the specified path
	rootURL := "https://raw.githubusercontent.com/deis/deis/"
	tag := args["--tag"].(string)
	units := []string{
		"deis-builder.service",
		"deis-cache.service",
		"deis-controller.service",
		"deis-database.service",
		"deis-logger.service",
		"deis-logspout.service",
		"deis-publisher.service",
		"deis-registry.service",
		"deis-router.service",
		"deis-store-admin.service",
		"deis-store-daemon.service",
		"deis-store-gateway.service",
		"deis-store-metadata.service",
		"deis-store-monitor.service",
		"deis-store-volume.service",
	}
	for _, unit := range units {
		src := rootURL + tag + "/deisctl/units/" + unit
		dest := filepath.Join(dir, unit)
		res, err := http.Get(src)
		if err != nil {
			return err
		}
		if res.StatusCode != 200 {
			return errors.New(res.Status)
		}
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(dest, data, 0644); err != nil {
			return err
		}
		fmt.Printf("Refreshed %s from %s\n", unit, tag)
	}
	return nil
}
