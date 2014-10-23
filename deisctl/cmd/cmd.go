package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/config"
	"github.com/deis/deis/deisctl/constant"
	"github.com/deis/deis/deisctl/update"
	"github.com/deis/deis/deisctl/utils"

	docopt "github.com/docopt/docopt-go"
)

const (
	PlatformInstallCommand string = "platform"
)

var (
	DefaultDataContainers = []string{
		"logger-data",
	}
)

func ListUnits(b backend.Backend) error {
	err := b.ListUnits()
	return err
}

func ListUnitFiles(b backend.Backend) error {
	err := b.ListUnitFiles()
	return err
}

func Scale(b backend.Backend, targets []string) error {
	for _, target := range targets {
		component, num, err := splitScaleTarget(target)
		if err != nil {
			return err
		}
		// the router is the only component that can scale past 1 at the moment
		if num > 1 && !strings.Contains(component, "router") {
			return fmt.Errorf("cannot scale %s past 1", component)
		}
		if err := b.Scale(component, num); err != nil {
			return err
		}
	}
	return nil
}

func Start(b backend.Backend, targets []string) error {

	if len(targets) == 1 && targets[0] == PlatformInstallCommand {
		return StartPlatform(b)
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

func StartPlatform(b backend.Backend) error {

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Starting Deis...")

	startDataContainers(b, &wg, outchan, errchan)
	startDefaultServices(b, &wg, outchan, errchan)

	wg.Wait()
	close(outchan)

	fmt.Println("Done.")
	fmt.Println()
	fmt.Println("Please use `deis register` to setup an administrator account.")
	return nil
}

func startDataContainers(b backend.Backend, wg *sync.WaitGroup, outchan chan string, errchan chan error) {
	outchan <- fmt.Sprintf("Data containers...")
	b.Start(DefaultDataContainers, wg, outchan, errchan)
	wg.Wait()
}

func startDefaultServices(b backend.Backend, wg *sync.WaitGroup, outchan chan string, errchan chan error) {

	// create separate channels for background tasks
	_outchan := make(chan string)
	_errchan := make(chan error)
	var _wg sync.WaitGroup

	outchan <- fmt.Sprintf("Logging subsystem...")
	b.Start([]string{"logger"}, wg, outchan, errchan)
	wg.Wait()
	b.Start([]string{"logspout"}, wg, outchan, errchan)
	wg.Wait()

	// start all services in the background
	b.Start([]string{"cache", "database", "registry", "controller", "builder",
		"publisher", "router@1", "router@2", "router@3"}, &_wg, _outchan, _errchan)

	// wait for groups to come up
	outchan <- fmt.Sprintf("Storage subsystem...")
	b.Start([]string{"store-daemon", "store-monitor", "store-gateway"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Control plane...")
	b.Start([]string{"cache", "database", "registry", "controller"}, wg, outchan, errchan)
	wg.Wait()

	b.Start([]string{"builder"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Data plane...")
	b.Start([]string{"publisher"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Routing mesh...")
	b.Start([]string{"router@1", "router@2", "router@3"}, wg, outchan, errchan)
	wg.Wait()
}

func Stop(b backend.Backend, targets []string) error {

	if len(targets) == 1 && targets[0] == PlatformInstallCommand {
		return StopPlatform(b)
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

func StopPlatform(b backend.Backend) error {

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Stopping Deis...")

	stopDefaultServices(b, &wg, outchan, errchan)

	wg.Wait()
	close(outchan)

	fmt.Println("Done.")
	fmt.Println()
	fmt.Println("Please run `deisctl start platform` to restart Deis.")
	return nil
}

func stopDefaultServices(b backend.Backend, wg *sync.WaitGroup, outchan chan string, errchan chan error) {

	outchan <- fmt.Sprintf("Routing mesh...")
	b.Stop([]string{"router@1", "router@2", "router@3"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Data plane...")
	b.Stop([]string{"publisher"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Control plane...")
	b.Stop([]string{"controller", "builder", "cache", "database", "registry"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Storage subsystem...")
	b.Stop([]string{"store-gateway", "store-monitor", "store-daemon"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Logging subsystem...")
	b.Stop([]string{"logger", "logspout"}, wg, outchan, errchan)
	wg.Wait()
}

func Restart(b backend.Backend, targets []string) error {
	if err := Stop(b, targets); err != nil {
		return err
	}
	return Start(b, targets)
}

func Status(b backend.Backend, targets []string) error {
	for _, target := range targets {
		if err := b.Status(target); err != nil {
			return err
		}
	}
	return nil
}

func Journal(b backend.Backend, targets []string) error {
	for _, target := range targets {
		if err := b.Journal(target); err != nil {
			return err
		}
	}
	return nil
}

func Install(b backend.Backend, targets []string) error {

	// if target is platform, install all services
	if len(targets) == 1 && targets[0] == PlatformInstallCommand {
		return InstallPlatform(b)
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

func InstallPlatform(b backend.Backend) error {

	if err := checkRequiredKeys(); err != nil {
		return err
	}

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Installing Deis...")

	installDataContainers(b, &wg, outchan, errchan)
	installDefaultServices(b, &wg, outchan, errchan)

	wg.Wait()
	close(outchan)

	fmt.Println("Done.")
	fmt.Println()
	fmt.Println("Please run `deisctl start platform` to boot up Deis.")
	return nil
}

func installDataContainers(b backend.Backend, wg *sync.WaitGroup, outchan chan string, errchan chan error) {
	outchan <- fmt.Sprintf("Data containers...")
	b.Create(DefaultDataContainers, wg, outchan, errchan)
	wg.Wait()
}

func installDefaultServices(b backend.Backend, wg *sync.WaitGroup, outchan chan string, errchan chan error) {

	outchan <- fmt.Sprintf("Logging subsystem...")
	b.Create([]string{"logger", "logspout"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Storage subsystem...")
	b.Create([]string{"store-daemon", "store-monitor", "store-gateway"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Control plane...")
	b.Create([]string{"cache", "database", "registry", "controller", "builder"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Data plane...")
	b.Create([]string{"publisher"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Routing mesh...")
	b.Create([]string{"router@1", "router@2", "router@3"}, wg, outchan, errchan)
	wg.Wait()
}

func Uninstall(b backend.Backend, targets []string) error {

	// if target is platform, uninstall all services
	if len(targets) == 1 && targets[0] == PlatformInstallCommand {
		return UninstallPlatform(b)
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

func UninstallPlatform(b backend.Backend) error {

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Uninstalling Deis...")

	uninstallAllServices(b, &wg, outchan, errchan)

	wg.Wait()
	close(outchan)

	fmt.Println("Done.")
	return nil
}

func uninstallAllServices(b backend.Backend, wg *sync.WaitGroup, outchan chan string, errchan chan error) error {

	outchan <- fmt.Sprintf("Routing mesh...")
	b.Destroy([]string{"router@1", "router@2", "router@3"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Data plane...")
	b.Destroy([]string{"publisher"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Control plane...")
	b.Destroy([]string{"controller", "builder", "cache", "database", "registry"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Storage subsystem...")
	b.Destroy([]string{"store-gateway", "store-monitor", "store-daemon"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Logging subsystem...")
	b.Destroy([]string{"logger", "logspout"}, wg, outchan, errchan)
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

func Config() error {
	if err := config.Config(); err != nil {
		return err
	}
	return nil
}

func Update() error {
	if err := utils.Execute(constant.HooksDir + "pre-update"); err != nil {
		fmt.Println("pre-updatehook failed")
		return err
	}
	if err := update.Update(); err != nil {
		fmt.Println("update engine failed")
		return err
	}
	if err := utils.Execute(constant.HooksDir + "post-update"); err != nil {
		fmt.Println("post-updatehook failed")
		return err
	}
	return nil
}

func RefreshUnits() error {
	usage := `Refreshes local unit files from the master repository.

deisctl looks for unit files in these directories, in this order:
- the $DEISCTL_UNITS environment variable, if set
- $HOME/.deis/units
- /var/lib/deis/units

Usage:
  deisctl refresh-units [-p <target>] [-t <tag>]

Options:
  -p --path=<target>   where to save unit files [default: $HOME/.deis/units]
  -t --tag=<tag>       git tag, branch, or SHA to use when downloading unit files
                       [default: v0.14.1]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, nil, true, "", false)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}
	dir := args["--path"].(string)
	if dir == "$HOME/.deis/units" || dir == "~/.deis/units" {
		dir = path.Join(os.Getenv("HOME"), ".deis", "units")
	}
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
		"deis-logger-data.service",
		"deis-logspout.service",
		"deis-publisher.service",
		"deis-registry.service",
		"deis-router.service",
		"deis-store-daemon.service",
		"deis-store-gateway.service",
		"deis-store-monitor.service",
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
