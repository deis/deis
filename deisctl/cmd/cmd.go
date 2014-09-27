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

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/config"
	"github.com/deis/deis/deisctl/constant"
	"github.com/deis/deis/deisctl/update"
	"github.com/deis/deis/deisctl/utils"
	"github.com/docopt/docopt-go"
)

const (
	PlatformInstallCommand string = "platform"
)

var (
	DefaultDataContainers = []string{
		"database-data",
		"registry-data",
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
	// if target is platform, start all services
	if len(targets) == 1 && targets[0] == PlatformInstallCommand {
		return StartPlatform(b)
	}
	return b.Start(targets)
}

func StartPlatform(b backend.Backend) error {
	fmt.Println(utils.DeisIfy("Starting Deis..."))
	if err := startDataContainers(b); err != nil {
		return err
	}
	if err := startDefaultServices(b); err != nil {
		return err
	}
	fmt.Println("Deis started.")
	return nil
}

func startDataContainers(b backend.Backend) error {
	fmt.Println("Launching data containers...")
	if err := b.Start(DefaultDataContainers); err != nil {
		return err
	}
	fmt.Println("Data containers launched.")
	return nil
}

func startDefaultServices(b backend.Backend) error {
	fmt.Println("Launching service containers...")
	if err := Start(b, []string{"logger"}); err != nil {
		return err
	}
	if err := Start(b, []string{"cache", "router", "database", "controller", "registry", "builder"}); err != nil {
		return err
	}
	fmt.Println("Service containers launched.")
	return nil
}

func Stop(b backend.Backend, targets []string) error {
	// if target is platform, stop all services
	if len(targets) == 1 && targets[0] == PlatformInstallCommand {
		return StopPlatform(b)
	}
	return b.Stop(targets)
}

func StopPlatform(b backend.Backend) error {
	fmt.Println("Stopping Deis...")
	if err := stopDefaultServices(b); err != nil {
		return err
	}
	fmt.Println("Deis stopped.")
	return nil
}

func stopDefaultServices(b backend.Backend) error {
	fmt.Println("Stopping service containers...")
	if err := Stop(b, []string{"builder", "registry", "controller", "database", "cache", "router", "logger"}); err != nil {
		return err
	}
	fmt.Println("Service containers stopped.")
	return nil
}

func Restart(b backend.Backend, targets []string) error {
	if err := b.Stop(targets); err != nil {
		return err
	}
	return b.Start(targets)
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
	// otherwise create the specific targets
	for i, target := range targets {
		// if we're installing a component without a number attached,
		// consider the user doesn't know better
		if !strings.Contains(target, "@") {
			targets[i] += "@1"
		}
	}
	return b.Create(targets)
}

func InstallPlatform(b backend.Backend) error {
	fmt.Println(utils.DeisIfy("Installing Deis..."))
	if err := installDataContainers(b); err != nil {
		return err
	}
	if err := installDefaultServices(b); err != nil {
		return err
	}
	fmt.Println("Deis installed.")
	fmt.Println("Please run `deisctl start platform` to boot up Deis.")
	return nil
}

func installDataContainers(b backend.Backend) error {
	fmt.Println("Scheduling data containers...")
	if err := b.Create(DefaultDataContainers); err != nil {
		return err
	}
	fmt.Println("Data containers scheduled.")
	return nil
}

func installDefaultServices(b backend.Backend) error {
	// start service containers
	targets := []string{
		"database=1",
		"cache=1",
		"logger=1",
		"registry=1",
		"controller=1",
		"builder=1",
		"router=1"}
	fmt.Println("Scheduling service containers...")
	if err := Scale(b, targets); err != nil {
		return err
	}
	fmt.Println("Service containers scheduled.")
	return nil
}

func Uninstall(b backend.Backend, targets []string) error {
	// if target is platform, uninstall all services
	if len(targets) == 1 && targets[0] == PlatformInstallCommand {
		return uninstallAllServices(b)
	}
	// uninstall the specific target
	return b.Destroy(targets)
}

func uninstallAllServices(b backend.Backend) error {
	targets := []string{
		"database=0",
		"cache=0",
		"logger=0",
		"registry=0",
		"controller=0",
		"builder=0",
		"router=0"}
	fmt.Println("Destroying service containers...")
	err := Scale(b, targets)
	fmt.Println("Service containers destroyed.")
	return err
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
  -p --path=<target>   where to save unit files [default: /var/lib/deis/units]
  -t --tag=<tag>       git tag, branch, or SHA to use when downloading unit files
                       [default: master]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, nil, true, "", false)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}
	dir := args["--path"].(string)
	// create the target dir if necessary
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	// download and save the unit files to the specified path
	rootURL := "https://raw.githubusercontent.com/deis/deis/deisctl/"
	tag := args["--tag"].(string)
	units := []string{
		"deis-builder.service",
		"deis-cache.service",
		"deis-controller.service",
		"deis-database.service",
		"deis-database-data.service",
		"deis-logger.service",
		"deis-logger-data.service",
		"deis-registry.service",
		"deis-registry-data.service",
		"deis-router.service",
	}
	for _, unit := range units {
		src := rootURL + tag + "/units/" + unit
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
