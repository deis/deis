package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/config"
	"github.com/deis/deis/deisctl/units"
	"github.com/deis/deis/deisctl/utils"
	"github.com/deis/deis/deisctl/utils/net"
)

const (
	// PlatformCommand is shorthand for "all the Deis components."
	PlatformCommand string = "platform"
	// StatelessPlatformCommand is shorthand for the components except store-*, database, and logger.
	StatelessPlatformCommand string = "stateless-platform"
	swarm                    string = "swarm"
	mesos                    string = "mesos"
	// DefaultRouterMeshSize defines the default number of routers to be loaded when installing the platform.
	DefaultRouterMeshSize uint8  = 3
	k8s                   string = "k8s"
)

// ListUnits prints a list of installed units.
func ListUnits(b backend.Backend) error {
	return b.ListUnits()
}

// ListUnitFiles prints the contents of all defined unit files.
func ListUnitFiles(b backend.Backend) error {
	return b.ListUnitFiles()
}

// Location to write standard output. By default, this is the os.Stdout.
var Stdout io.Writer = os.Stderr

// Location to write standard error information. By default, this is the os.Stderr.
var Stderr io.Writer = os.Stdout

// Number of routers to be installed. By default, it's DefaultRouterMeshSize.
var RouterMeshSize = DefaultRouterMeshSize

// Scale grows or shrinks the number of running components.
// Currently "router", "registry" and "store-gateway" are the only types that can be scaled.
func Scale(targets []string, b backend.Backend) error {
	var wg sync.WaitGroup

	for _, target := range targets {
		component, num, err := splitScaleTarget(target)
		if err != nil {
			return err
		}
		// the router, registry, and store-gateway are the only component that can scale at the moment
		if !strings.Contains(component, "router") && !strings.Contains(component, "registry") && !strings.Contains(component, "store-gateway") {
			return fmt.Errorf("cannot scale %s component", component)
		}
		b.Scale(component, num, &wg, Stdout, Stderr)
		wg.Wait()
	}
	return nil
}

// Start activates the specified components.
func Start(targets []string, b backend.Backend) error {

	// if target is platform, install all services
	if len(targets) == 1 {
		switch targets[0] {
		case PlatformCommand:
			return StartPlatform(b, false)
		case StatelessPlatformCommand:
			return StartPlatform(b, true)
		case mesos:
			return StartMesos(b)
		case swarm:
			return StartSwarm(b)
		case k8s:
			return StartK8s(b)
		}
	}
	var wg sync.WaitGroup

	b.Start(targets, &wg, Stdout, Stderr)
	wg.Wait()

	return nil
}

// RollingRestart restart instance unit in a rolling manner
func RollingRestart(target string, b backend.Backend) error {
	var wg sync.WaitGroup

	b.RollingRestart(target, &wg, Stdout, Stderr)
	wg.Wait()

	return nil
}

// CheckRequiredKeys exist in config backend
func CheckRequiredKeys(cb config.Backend) error {
	if err := config.CheckConfig("/deis/platform/", "domain", cb); err != nil {
		return fmt.Errorf(`Missing platform domain, use:
deisctl config platform set domain=<your-domain>`)
	}

	if err := config.CheckConfig("/deis/platform/", "sshPrivateKey", cb); err != nil {
		fmt.Printf(`Warning: Missing sshPrivateKey, "deis run" will be unavailable. Use:
deisctl config platform set sshPrivateKey=<path-to-key>
`)
	}
	return nil
}

func startDefaultServices(b backend.Backend, stateless bool, wg *sync.WaitGroup, out, err io.Writer) {

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
		"publisher", "router@*",
	}
	if stateless {
		batch = []string{"registry@*", "controller", "builder", "publisher", "router@*"}
	}
	b.Start(batch, &bgwg, &trash, &trash)
	// End background stuff.

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

	fmt.Fprintln(out, "Router mesh...")
	b.Start([]string{"router@*"}, wg, out, err)
	wg.Wait()
}

// Stop deactivates the specified components.
func Stop(targets []string, b backend.Backend) error {

	// if target is platform, stop all services
	if len(targets) == 1 {
		switch targets[0] {
		case PlatformCommand:
			return StopPlatform(b, false)
		case StatelessPlatformCommand:
			return StopPlatform(b, true)
		case mesos:
			return StopMesos(b)
		case swarm:
			return StopSwarm(b)
		case k8s:
			return StopK8s(b)
		}
	}

	var wg sync.WaitGroup

	b.Stop(targets, &wg, Stdout, Stderr)
	wg.Wait()

	return nil
}

func stopDefaultServices(b backend.Backend, stateless bool, wg *sync.WaitGroup, out, err io.Writer) {

	fmt.Fprintln(out, "Router mesh...")
	b.Stop([]string{"router@*"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Data plane...")
	b.Stop([]string{"publisher"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Control plane...")
	if stateless {
		b.Stop([]string{"controller", "builder", "registry@*"}, wg, out, err)
	} else {
		b.Stop([]string{"controller", "builder", "database", "registry@*"}, wg, out, err)
	}
	wg.Wait()

	fmt.Fprintln(out, "Logging subsystem...")
	if stateless {
		b.Stop([]string{"logspout"}, wg, out, err)
	} else {
		b.Stop([]string{"logger", "logspout"}, wg, out, err)
	}
	wg.Wait()

	if !stateless {
		fmt.Fprintln(out, "Storage subsystem...")
		b.Stop([]string{"store-volume", "store-gateway@*"}, wg, out, err)
		wg.Wait()
		b.Stop([]string{"store-metadata"}, wg, out, err)
		wg.Wait()
		b.Stop([]string{"store-daemon"}, wg, out, err)
		wg.Wait()
		b.Stop([]string{"store-monitor"}, wg, out, err)
		wg.Wait()
	}

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
func Install(targets []string, b backend.Backend, cb config.Backend, checkKeys func(config.Backend) error) error {

	// if target is platform, install all services
	if len(targets) == 1 {
		switch targets[0] {
		case PlatformCommand:
			return InstallPlatform(b, cb, checkKeys, false)
		case StatelessPlatformCommand:
			return InstallPlatform(b, cb, checkKeys, true)
		case mesos:
			return InstallMesos(b)
		case swarm:
			return InstallSwarm(b)
		case k8s:
			return InstallK8s(b)
		}
	}
	var wg sync.WaitGroup

	// otherwise create the specific targets
	b.Create(targets, &wg, Stdout, Stderr)
	wg.Wait()

	return nil
}

func installDefaultServices(b backend.Backend, stateless bool, wg *sync.WaitGroup, out, err io.Writer) {

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

	fmt.Fprintln(out, "Router mesh...")
	b.Create(getRouters(), wg, out, err)
	wg.Wait()

}

func getRouters() []string {
	routers := make([]string, RouterMeshSize)
	for i := uint8(0); i < RouterMeshSize; i++ {
		routers[i] = fmt.Sprintf("router@%d", i+1)
	}
	return routers
}

// Uninstall unloads the definitions of the specified components.
// After Uninstall, the components will be unavailable until Install is called.
func Uninstall(targets []string, b backend.Backend) error {
	if len(targets) == 1 {
		switch targets[0] {
		case PlatformCommand:
			return UninstallPlatform(b, false)
		case StatelessPlatformCommand:
			return UninstallPlatform(b, true)
		case mesos:
			return UninstallMesos(b)
		case swarm:
			return UnInstallSwarm(b)
		case k8s:
			return UnInstallK8s(b)
		}
	}

	var wg sync.WaitGroup

	// uninstall the specific target
	b.Destroy(targets, &wg, Stdout, Stderr)
	wg.Wait()

	return nil
}

func uninstallAllServices(b backend.Backend, stateless bool, wg *sync.WaitGroup, out, err io.Writer) error {

	fmt.Fprintln(out, "Router mesh...")
	b.Destroy([]string{"router@*"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Data plane...")
	b.Destroy([]string{"publisher"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Control plane...")
	if stateless {
		b.Destroy([]string{"controller", "builder", "registry@*"}, wg, out, err)
	} else {
		b.Destroy([]string{"controller", "builder", "database", "registry@*"}, wg, out, err)
	}
	wg.Wait()

	fmt.Fprintln(out, "Logging subsystem...")
	if stateless {
		b.Destroy([]string{"logspout"}, wg, out, err)
	} else {
		b.Destroy([]string{"logger", "logspout"}, wg, out, err)
	}
	wg.Wait()

	if !stateless {
		fmt.Fprintln(out, "Storage subsystem...")
		b.Destroy([]string{"store-volume", "store-gateway@*"}, wg, out, err)
		wg.Wait()
		b.Destroy([]string{"store-metadata"}, wg, out, err)
		wg.Wait()
		b.Destroy([]string{"store-daemon"}, wg, out, err)
		wg.Wait()
		b.Destroy([]string{"store-monitor"}, wg, out, err)
		wg.Wait()
	}

	return nil
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
// A configuration value is stored and retrieved from a key/value store
// at /deis/<component>/<config>. Configuration values are typically used for component-level
// configuration, such as enabling TLS for the routers.
func Config(target string, action string, key []string, cb config.Backend) error {
	if err := config.Config(target, action, key, cb); err != nil {
		return err
	}
	return nil
}

// RefreshUnits overwrites local unit files with those requested.
// Downloading from the Deis project GitHub URL by tag or SHA is the only mechanism
// currently supported.
func RefreshUnits(unitDir, tag, rootURL string) error {
	unitDir = utils.ResolvePath(unitDir)
	decoratorDir := filepath.Join(unitDir, "decorators")
	// create the target dir if necessary
	if err := os.MkdirAll(decoratorDir, 0755); err != nil {
		return err
	}
	// download and save the unit files to the specified path
	for _, unit := range units.Names {
		unitSrc := rootURL + tag + "/deisctl/units/" + unit + ".service"
		unitDest := filepath.Join(unitDir, unit+".service")
		if err := net.Download(unitSrc, unitDest); err != nil {
			return err
		}
		fmt.Printf("Refreshed %s unit from %s\n", unit, tag)
		decoratorSrc := rootURL + tag + "/deisctl/units/decorators/" + unit + ".service.decorator"
		decoratorDest := filepath.Join(decoratorDir, unit+".service.decorator")
		if err := net.Download(decoratorSrc, decoratorDest); err != nil {
			if err.Error() == "404 Not Found" {
				fmt.Printf("Decorator for %s not found in %s\n", unit, tag)
			} else {
				return err
			}
		} else {
			fmt.Printf("Refreshed %s decorator from %s\n", unit, tag)
		}
	}
	return nil
}

// SSH opens an interactive shell on a machine in the cluster
func SSH(target string, cmd []string, b backend.Backend) error {

	if len(cmd) > 0 {
		return b.SSHExec(target, strings.Join(cmd, " "))
	}

	return b.SSH(target)
}

// Dock connects to the appropriate host and runs 'docker exec -it'.
func Dock(target string, cmd []string, b backend.Backend) error {
	return b.Dock(target, cmd)
}
