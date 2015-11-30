package cmd

import (
	"fmt"
	"io"
	"sync"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/config"
	"github.com/deis/deis/pkg/prettyprint"
)

// InstallPlatform loads all components' definitions from local unit files.
// After InstallPlatform, all components will be available for StartPlatform.
func InstallPlatform(b backend.Backend, cb config.Backend, checkKeys func(config.Backend) error, stateless bool) error {

	if err := checkKeys(cb); err != nil {
		return err
	}

	if stateless {
		fmt.Println("Warning: With a stateless control plane, some components require manual configuration.")
		fmt.Println("See the official Deis documentation for details on running a stateless control plane.")
		fmt.Println("http://docs.deis.io/en/latest/managing_deis/running-deis-without-ceph/")
	}

	var wg sync.WaitGroup

	io.WriteString(Stdout, prettyprint.DeisIfy("Installing Deis..."))

	installDefaultServices(b, stateless, &wg, Stdout, Stderr)

	wg.Wait()

	fmt.Fprintln(Stdout, "Done.")
	fmt.Fprintln(Stdout, "")
	if stateless {
		fmt.Fprintln(Stdout, "Please run `deisctl start stateless-platform` to boot up Deis.")
	} else {
		fmt.Fprintln(Stdout, "Please run `deisctl start platform` to boot up Deis.")
	}
	return nil
}

// StartPlatform activates all components.
func StartPlatform(b backend.Backend, stateless bool) error {

	var wg sync.WaitGroup

	io.WriteString(Stdout, prettyprint.DeisIfy("Starting Deis..."))

	startDefaultServices(b, stateless, &wg, Stdout, Stderr)

	wg.Wait()

	fmt.Fprintln(Stdout, "Done.\n ")
	fmt.Fprintln(Stdout, "Please set up an administrative account. See 'deis help register'")
	return nil
}

// StopPlatform deactivates all components.
func StopPlatform(b backend.Backend, stateless bool) error {

	var wg sync.WaitGroup

	io.WriteString(Stdout, prettyprint.DeisIfy("Stopping Deis..."))

	stopDefaultServices(b, stateless, &wg, Stdout, Stderr)

	wg.Wait()

	fmt.Fprintln(Stdout, "Done.\n ")
	if stateless {
		fmt.Fprintln(Stdout, "Please run `deisctl start stateless-platform` to restart Deis.")
	} else {
		fmt.Fprintln(Stdout, "Please run `deisctl start platform` to restart Deis.")
	}
	return nil
}

// UninstallPlatform unloads all components' definitions.
// After UninstallPlatform, all components will be unavailable.
func UninstallPlatform(b backend.Backend, stateless bool) error {

	var wg sync.WaitGroup

	io.WriteString(Stdout, prettyprint.DeisIfy("Uninstalling Deis..."))

	uninstallAllServices(b, stateless, &wg, Stdout, Stderr)

	wg.Wait()

	fmt.Fprintln(Stdout, "Done.")
	return nil
}
