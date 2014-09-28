package cmd

import (
	"fmt"

	"github.com/deis/deisctl/backend"
	"github.com/deis/deisctl/utils"
)

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
