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
	"github.com/deis/deis/deisctl/units"
	"github.com/deis/deis/deisctl/utils"
)

const (
	// PlatformCommand is shorthand for "all the Deis components."
	PlatformCommand string = "platform"
	swarm           string = "swarm"
)

// ListUnits prints a list of installed units.
func ListUnits(b backend.Backend) error {
	return b.ListUnits()
}

// ListUnitFiles prints the contents of all defined unit files.
func ListUnitFiles(b backend.Backend) error {
	return b.ListUnitFiles()
}

// Scale grows or shrinks the number of running components.
// Currently "router", "registry" and "store-gateway" are the only types that can be scaled.
func Scale(targets []string, b backend.Backend) error {
	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	for _, target := range targets {
		component, num, err := splitScaleTarget(target)
		if err != nil {
			return err
		}
		// the router, registry, and store-gateway are the only component that can scale at the moment
		if !strings.Contains(component, "router") && !strings.Contains(component, "registry") && !strings.Contains(component, "store-gateway") {
			return fmt.Errorf("cannot scale %s component", component)
		}
		b.Scale(component, num, &wg, outchan, errchan)
		wg.Wait()
	}
	close(outchan)
	close(errchan)
	return nil
}

// Start activates the specified components.
func Start(targets []string, b backend.Backend) error {

	// if target is platform, install all services
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
	close(errchan)

	return nil
}

// CheckRequiredKeys exist in etcd
func CheckRequiredKeys() error {
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
func Stop(targets []string, b backend.Backend) error {

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
	close(errchan)

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
func Restart(targets []string, b backend.Backend) error {

	// act as if the user called "stop" and then "start"
	if err := Stop(targets, b); err != nil {
		return err
	}

	return Start(targets, b)
}

// Status prints the current status of components.
func Status(targets []string, b backend.Backend) error {

	for _, target := range targets {
		if err := b.Status(target); err != nil {
			return err
		}
	}
	return nil
}

// Journal prints log output for the specified components.
func Journal(targets []string, b backend.Backend) error {

	for _, target := range targets {
		if err := b.Journal(target); err != nil {
			return err
		}
	}
	return nil
}

// Install loads the definitions of components from local unit files.
// After Install, the components will be available to Start.
func Install(targets []string, b backend.Backend, checkKeys func() error) error {

	// if target is platform, install all services
	if len(targets) == 1 {
		if targets[0] == PlatformCommand {
			return InstallPlatform(b, checkKeys)
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
	close(errchan)
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
func Uninstall(targets []string, b backend.Backend) error {
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
	close(errchan)

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

func printState(outchan chan string, errchan chan error, interval time.Duration) {
	for {
		select {
		case out, ok := <-outchan:
			if !ok {
				outchan = nil
			}
			if out != "" {
				fmt.Println(out)
			}
		case err, ok := <-errchan:
			if !ok {
				errchan = nil
			}
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		if outchan == nil && errchan == nil {
			break
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
func Config(target string, action string, key []string) error {
	if err := config.Config(target, action, key); err != nil {
		return err
	}
	return nil
}

// RefreshUnits overwrites local unit files with those requested.
// Downloading from the Deis project GitHub URL by tag or SHA is the only mechanism
// currently supported.
func RefreshUnits(dir, tag, url string) error {
	dir = utils.ResolvePath(dir)
	// create the target dir if necessary
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	// download and save the unit files to the specified path
	for _, unit := range units.Names {
		src := fmt.Sprintf(url, tag, unit)
		dest := filepath.Join(dir, unit+".service")
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

// SSH opens an interactive shell on a machine in the cluster
func SSH(target string, b backend.Backend) error {
	if err := b.SSH(target); err != nil {
		return err
	}

	return nil
}
