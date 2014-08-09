package cmd

import (
	"fmt"
	"github.com/deis/deisctl/client"
	"github.com/deis/deisctl/utils"
	"regexp"
	"strconv"
	"strings"
)

func List(c client.Client) error {
	err := c.List()
	return err
}

func PullImage(service string) error {
	dockercli, _, _ := utils.GetNewClient()
	fmt.Println("pulling image :" + strings.Split(service, ".")[0])
	err := utils.CmdPull(dockercli, strings.Split(service, ".")[0])
	if err != nil {
		return err
	}
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

func Status(c client.Client, targets []string) error {
	for _, target := range targets {
		err := c.Status(target)
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

func Update(args []string) {

	if len(args) != 4 {
		fmt.Println("unsufficient args")
		fmt.Println("usage:  updatectl update instance deis")
		return
	}
	if args[2] != "instance" && args[3] != "deis" {
		fmt.Println("wrong args ")
		fmt.Println("usage:  updatectl update instance deis")
		return
	}
	Args := []string{
		"instance",
		"deis",
		"--clients-per-app=1",
		"--min-sleep=5",
		"--max-sleep=10",
		"--app-id=329cd607-06fe-4bde-8ecd-613b58c6945f",
		"--group-id=bee2027e-29a4-4135-bffb-b2864234dd15",
		"--version=1.1.0",
	}
	updatectl.Update(Args)
}

func installDataContainers(c client.Client) error {
	// data containers
	dataContainers := []string{
		"database-data",
		"registry-data",
		"logger-data",
		"builder-data",
	}
	fmt.Println("Scheduling data containers...")
	for _, dataContainer := range dataContainers {
		c.Create(dataContainer)
		// if err != nil {
		// 	return err
		// }
	}
	fmt.Println("Activating data containers...")
	for _, dataContainer := range dataContainers {
		c.Start(dataContainer)
		// if err != nil {
		// 	return err
		// }
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
	fmt.Println("Scheduling units...")
	err := Scale(c, targets)
	fmt.Println("Activating units...")
	err = Start(c, []string{"registry", "logger", "cache", "database"})
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
	// if targets, uninstall all services
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
	err := Scale(c, targets)
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
