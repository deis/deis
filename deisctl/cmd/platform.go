package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/utils"
)

// InstallPlatform loads all components' definitions from local unit files.
// After InstallPlatform, all components will be available for StartPlatform.
func InstallPlatform(b backend.Backend) error {

	if err := checkRequiredKeys(); err != nil {
		return err
	}

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Installing Deis...")

	installDefaultServices(b, &wg, outchan, errchan)

	wg.Wait()
	close(outchan)

	fmt.Println("Done.")
	fmt.Println()
	fmt.Println("Please run `deisctl start platform` to boot up Deis.")
	return nil
}

// StartPlatform activates all components.
func StartPlatform(b backend.Backend) error {

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Starting Deis...")

	startDefaultServices(b, &wg, outchan, errchan)

	wg.Wait()
	close(outchan)

	fmt.Println("Done.")
	fmt.Println()
	fmt.Println("Please use `deis register` to setup an administrator account.")
	return nil
}

// StopPlatform deactivates all components.
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

// UninstallPlatform unloads all components' definitions.
// After UninstallPlatform, all components will be unavailable.
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
