package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/deis/deisctl/client"
	"github.com/deis/deisctl/config"
	"github.com/deis/deisctl/constant"
	"github.com/deis/deisctl/update"
	"github.com/deis/deisctl/utils"
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
		"builder-data",
	}
)

func ListUnits(c client.Client) error {
	err := c.ListUnits()
	return err
}

func ListUnitFiles(c client.Client) error {
	err := c.ListUnitFiles()
	return err
}

func Scale(c client.Client, targets []string) error {
	for _, target := range targets {
		component, num, err := splitScaleTarget(target)
		if err != nil {
			return err
		}
		err = c.Scale(component, num)
		if err != nil {
			return err
		}
	}
	return nil
}

func Start(c client.Client, targets []string) error {
	// if target is platform, start all services
	if len(targets) == 1 && targets[0] == PlatformInstallCommand {
		return StartPlatform(c)
	}
	return c.Start(targets)
}

func StartPlatform(c client.Client) error {
	fmt.Println("Starting Platform...")
	if err := startDataContainers(c); err != nil {
		return err
	}
	if err := startDefaultServices(c); err != nil {
		return err
	}
	fmt.Println("Platform started.")
	return nil
}

func startDataContainers(c client.Client) error {
	fmt.Println("Launching data containers...")
	if err := c.Start(DefaultDataContainers); err != nil {
		return err
	}
	fmt.Println("Data containers launched.")
	return nil
}

func startDefaultServices(c client.Client) error {
	fmt.Println("Launching service containers...")
	if err := Start(c, []string{"logger"}); err != nil {
		return err
	}
	if err := Start(c, []string{"cache", "router", "database", "controller", "registry", "builder"}); err != nil {
		return err
	}
	fmt.Println("Service containers launched.")
	return nil
}

func Stop(c client.Client, targets []string) error {
	// if target is platform, stop all services
	if len(targets) == 1 && targets[0] == PlatformInstallCommand {
		return StopPlatform(c)
	}
	return c.Stop(targets)
}

func StopPlatform(c client.Client) error {
	fmt.Println("Stopping Platform...")
	if err := stopDefaultServices(c); err != nil {
		return err
	}
	fmt.Println("Platform stopped.")
	return nil
}

func stopDefaultServices(c client.Client) error {
	fmt.Println("Stopping service containers...")
	if err := Stop(c, []string{"builder", "registry", "controller", "database", "cache", "router", "logger"}); err != nil {
		return err
	}
	fmt.Println("Service containers stopped.")
	return nil
}

func Restart(c client.Client, targets []string) error {
	if err := c.Stop(targets); err != nil {
		return err
	}
	return c.Start(targets)
}

func Status(c client.Client, targets []string) error {
	for _, target := range targets {
		if err := c.Status(target); err != nil {
			return err
		}
	}
	return nil
}

func Journal(c client.Client, targets []string) error {
	for _, target := range targets {
		if err := c.Journal(target); err != nil {
			return err
		}
	}
	return nil
}

func Install(c client.Client, targets []string) error {
	// if target is platform, install all services
	if len(targets) == 1 && targets[0] == PlatformInstallCommand {
		return InstallPlatform(c)
	}
	// otherwise create the specific targets
	for i, target := range targets {
		// if we're installing a component without a number attached,
		// consider the user doesn't know better
		if !strings.Contains(target, "@") {
			targets[i] += "@1"
		}
	}
	return c.Create(targets)
}

func InstallPlatform(c client.Client) error {
	if err := installDataContainers(c); err != nil {
		return err
	}
	return installDefaultServices(c)
}

func installDataContainers(c client.Client) error {
	fmt.Println("Scheduling data containers...")
	if err := c.Create(DefaultDataContainers); err != nil {
		return err
	}
	fmt.Println("Data containers scheduled.")
	return nil
}

func installDefaultServices(c client.Client) error {
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
	if err := Scale(c, targets); err != nil {
		return err
	}
	fmt.Println("Service containers scheduled.")
	return nil
}

func Uninstall(c client.Client, targets []string) error {
	// if target is platform, uninstall all services
	if len(targets) == 1 && targets[0] == PlatformInstallCommand {
		return uninstallAllServices(c)
	}
	// uninstall the specific target
	return c.Destroy(targets)
}

func uninstallAllServices(c client.Client) error {
	targets := []string{
		"database=0",
		"cache=0",
		"logger=0",
		"registry=0",
		"controller=0",
		"builder=0",
		"router=0"}
	fmt.Println("Destroying service containers...")
	err := Scale(c, targets)
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
  deisctl refresh-units [-p <target>]

Options:
  -p --path=<target>   where to save unit files [default: /var/lib/deis/units]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, nil, true, "", false)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}
	dir, _ := client.ExpandUser(args["--path"].(string))
	// create the target dir if necessary
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	// download and save the unit files to the specified path
	rootURL := "https://raw.githubusercontent.com/deis/deisctl/"
	branch := "master"
	units := []string{
		"deis-builder.service",
		"deis-builder-data.service",
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
		src := rootURL + branch + "/units/" + unit
		dest := filepath.Join(dir, unit)
		res, err := http.Get(src)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(dest, data, 0644); err != nil {
			return err
		}
		fmt.Printf("Refreshed %s from %s\n", unit, branch)
	}
	return nil
}
