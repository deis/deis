package cmd

import (
	"fmt"
	"strings"

	"github.com/deis/deisctl/backend"
	"github.com/deis/deisctl/constant"
	"github.com/deis/deisctl/update"
	"github.com/deis/deisctl/utils"
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
