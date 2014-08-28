package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/deis/deisctl/client"
	"github.com/deis/deisctl/config"
	"github.com/deis/deisctl/constant"
	"github.com/deis/deisctl/update"
	"github.com/deis/deisctl/utils"
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
	for _, target := range targets {
		err := c.Start(target)
		if err != nil {
			return err
		}
	}
	return nil
}

func Stop(c client.Client, targets []string) error {
	for _, target := range targets {
		err := c.Stop(target)
		if err != nil {
			return err
		}
	}
	return nil
}

func Restart(c client.Client, targets []string) error {
	for _, target := range targets {
		err := c.Stop(target)
		if err != nil {
			return err
		}
		err = c.Start(target)
		if err != nil {
			return err
		}
	}
	return nil
}

func Status(c client.Client, targets []string) error {
	for _, target := range targets {
		err := c.Status(target)
		if err != nil {
			return err
		}
	}
	return nil
}

func Journal(c client.Client, targets []string) error {
	for _, target := range targets {
		err := c.Journal(target)
		if err != nil {
			return err
		}
	}
	return nil
}

func Install(c client.Client, targets []string) error {
	// if targets, install all services
	if len(targets) == 0 {
		err := installDataContainers(c)
		if err != nil {
			return err
		}
		err = installDefaultServices(c)
		if err != nil {
			return err
		}
	} else {
		// otherwise create and start the specific targets
		for _, target := range targets {
			err := c.Create(target)
			if err != nil {
				return err
			}
			err = c.Start(target)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func installDataContainers(c client.Client) error {
	// data containers
	dataContainers := []string{
		"database-data",
		"registry-data",
		"logger-data",
		"builder-data",
	}
	fmt.Println("\nScheduling data containers...")
	for _, dataContainer := range dataContainers {
		err := c.Create(dataContainer)
		if err != nil {
			return err
		}
	}
	fmt.Println("\nLaunching data containers...")
	for _, dataContainer := range dataContainers {
		err := c.Start(dataContainer)
		if err != nil {
			return err
		}
	}
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
	fmt.Println("\nScheduling service containers...")
	err := Scale(c, targets)
	fmt.Println("\nLaunching service containers...")
	err = Start(c, []string{"logger", "cache", "database"})
	if err != nil {
		return err
	}
	err = Start(c, []string{"registry"})
	if err != nil {
		return err
	}
	err = Start(c, []string{"controller"})
	if err != nil {
		return err
	}
	err = Start(c, []string{"builder"})
	if err != nil {
		return err
	}
	err = Start(c, []string{"router"})
	if err != nil {
		return err
	}
	fmt.Println("Done.")
	return nil
}

func Uninstall(c client.Client, targets []string) error {
	// if no targets, uninstall all services
	if len(targets) == 0 {
		err := uninstallAllServices(c)
		if err != nil {
			return err
		}
	} else {
		// uninstall the specific target
		for _, target := range targets {
			err := c.Destroy(target)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
	fmt.Println("\nDestroying service containers...")
	err := Scale(c, targets)
	fmt.Println("Done.")
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
	// create the $HOME/.deisctl directory if necessary
	user, err := user.Current()
	if err != nil {
		return err
	}
	dir := filepath.Join(user.HomeDir, ".deisctl")
	if err = os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// download and save the unit files to $HOME/.deisctl
	rootUrl := "https://raw.githubusercontent.com/deis/deisctl/"
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
		src := rootUrl + branch + "/units/" + unit
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
		if err = ioutil.WriteFile(dest, data, 0600); err != nil {
			return err
		}
		fmt.Printf("Refreshed %s from %s\n", unit, branch)
	}
	return nil
}
