package cmd

import (
	"fmt"
	"io"
	"sync"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/pkg/prettyprint"
)

//InstallSwarm Installs swarm
func InstallSwarm(b backend.Backend) error {
	var wg sync.WaitGroup
	io.WriteString(Stdout, prettyprint.DeisIfy("Installing Swarm..."))
	fmt.Fprintln(Stdout, "Swarm node and Swarm Manager...")
	b.Create([]string{"swarm-node", "swarm-manager"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "Done.\n ")
	fmt.Fprintln(Stdout, "Please run `deisctl start swarm` to start swarm.")
	return nil
}

//StartSwarm starts Swarm Schduler
func StartSwarm(b backend.Backend) error {
	var wg sync.WaitGroup
	io.WriteString(Stdout, prettyprint.DeisIfy("Starting Swarm..."))
	fmt.Fprintln(Stdout, "swarm nodes...")
	b.Start([]string{"swarm-node"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "swarm manager...")
	b.Start([]string{"swarm-manager"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "Done.\n ")
	fmt.Fprintln(Stdout, "Please run `deisctl config controller set schedulerModule=swarm` to use the swarm scheduler.")
	return nil
}

//StopSwarm stops swarm
func StopSwarm(b backend.Backend) error {

	var wg sync.WaitGroup
	io.WriteString(Stdout, prettyprint.DeisIfy("Stopping Swarm..."))
	fmt.Fprintln(Stdout, "swarm nodes and swarm manager")
	b.Stop([]string{"swarm-node", "swarm-manager"}, &wg, Stdout, Stderr)
	wg.Wait()

	fmt.Fprintln(Stdout, "Done.\n ")
	return nil
}

//UnInstallSwarm uninstall Swarm
func UnInstallSwarm(b backend.Backend) error {
	var wg sync.WaitGroup
	io.WriteString(Stdout, prettyprint.DeisIfy("Destroying Swarm..."))
	fmt.Fprintln(Stdout, "swarm nodes and swarm manager...")
	b.Destroy([]string{"swarm-node", "swarm-manager"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "Done.\n ")
	return nil
}
