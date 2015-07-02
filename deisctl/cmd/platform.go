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
func InstallPlatform(b backend.Backend, checkKeys func() error, stateless bool) error {

	if err := checkKeys(); err != nil {
		return err
	}

	if stateless {
		fmt.Println("Warning: With a stateless control plane, `deis logs` will be unavailable.")
		fmt.Println("Additionally, components will need to be configured to use external persistent stores.")
		fmt.Println("See the official Deis documentation for details on running a stateless control plane.")
	}

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Installing Deis...")

	installDefaultServices(b, stateless, &wg, outchan, errchan)

	wg.Wait()
	close(outchan)

	fmt.Println("Done.")
	fmt.Println()
	if stateless {
		fmt.Println("Please run `deisctl start stateless-platform` to boot up Deis.")
	} else {
		fmt.Println("Please run `deisctl start platform` to boot up Deis.")
	}
	return nil
}

// StartPlatform activates all components.
func StartPlatform(b backend.Backend, stateless bool) error {

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Starting Deis...")

	startDefaultServices(b, stateless, &wg, outchan, errchan)

	wg.Wait()
	close(outchan)

	fmt.Println("Done.")
	fmt.Println()
	fmt.Println("Please use `deis register` to setup an administrator account.")
	return nil
}

// StopPlatform deactivates all components.
func StopPlatform(b backend.Backend, stateless bool) error {

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Stopping Deis...")

	stopDefaultServices(b, stateless, &wg, outchan, errchan)

	wg.Wait()
	close(outchan)

	fmt.Println("Done.")
	fmt.Println()
	if stateless {
		fmt.Println("Please run `deisctl start stateless-platform` to restart Deis.")
	} else {
		fmt.Println("Please run `deisctl start platform` to restart Deis.")
	}
	return nil
}

// UninstallPlatform unloads all components' definitions.
// After UninstallPlatform, all components will be unavailable.
func UninstallPlatform(b backend.Backend, stateless bool) error {

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Uninstalling Deis...")

	uninstallAllServices(b, stateless, &wg, outchan, errchan)

	wg.Wait()
	close(outchan)

	fmt.Println("Done.")
	return nil
}
