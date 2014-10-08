package cmd

import (
	"fmt"

	"regexp"
	"strconv"
	"strings"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/config"
	"github.com/deis/deis/deisctl/constant"
	"github.com/deis/deis/deisctl/update"
	"github.com/deis/deis/deisctl/utils"
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
	if err := Start(b, []string{"logger@1"}); err != nil {
		return err
	}
	targets := []string{
		"publisher",
		"store-monitor",
		"store-daemon",
		"store-gateway@1",
		"logspout",
		"cache@1",
		"router@1",
		"database@1",
		"controller@1",
		"registry@1",
		"builder@1",
	}
	if err := Start(b, targets); err != nil {
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
	targets := []string{
		"publisher",
		"logspout",
		"builder@1",
		"registry@1",
		"controller@1",
		"database@1",
		"store-gateway@1",
		"store-daemon",
		"store-monitor",
		"cache@1",
		"router@1",
		"logger@1",
	}
	if err := Stop(b, targets); err != nil {
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
	// Install global units
	if err := b.Create([]string{"publisher", "logspout", "store-monitor", "store-daemon"}); err != nil {
		return err
	}
	// start service containers
	targets := []string{
		"store-gateway=1",
		"database=1",
		"cache=1",
		"logger=1",
		"registry=1",
		"controller=1",
		"builder=1",
		"router=1",
	}
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
		"store-gateway=0",
		"database=0",
		"cache=0",
		"logger=0",
		"registry=0",
		"controller=0",
		"builder=0",
		"router=0",
	}
	fmt.Println("Destroying service containers...")
	if err := Scale(b, targets); err != nil {
		return err
	}
	// Uninstall global units
	if err := b.Destroy([]string{"publisher", "logspout", "store-monitor", "store-daemon"}); err != nil {
		return err
	}
	fmt.Println("Service containers destroyed.")
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
