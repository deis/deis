package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/utils"
)

//InstallSwarm Installs swarm
func InstallSwarm(b backend.Backend) error {
	outchan := make(chan string)
	errchan := make(chan error)
	defer close(outchan)
	defer close(errchan)
	var wg sync.WaitGroup
	go printState(outchan, errchan, 500*time.Millisecond)
	outchan <- utils.DeisIfy("Installing Swarm...")
	outchan <- fmt.Sprintf("Swarm node and Swarm Manager...")
	b.Create([]string{"swarm-node", "swarm-manager"}, &wg, outchan, errchan)
	wg.Wait()
	fmt.Println("Done.")
	fmt.Println()
	fmt.Println("Please run `deisctl start swarm` to start swarm.")
	return nil
}

//StartSwarm starts Swarm Schduler
func StartSwarm(b backend.Backend) error {
	outchan := make(chan string)
	errchan := make(chan error)
	defer close(outchan)
	defer close(errchan)
	var wg sync.WaitGroup
	go printState(outchan, errchan, 500*time.Millisecond)
	outchan <- utils.DeisIfy("Starting Swarm...")
	outchan <- fmt.Sprintf("swarm nodes...")
	b.Start([]string{"swarm-node"}, &wg, outchan, errchan)
	wg.Wait()
	outchan <- fmt.Sprintf("swarm manager...")
	b.Start([]string{"swarm-manager"}, &wg, outchan, errchan)
	wg.Wait()
	fmt.Println("Done.")
	fmt.Println("Please run `deisctl config controller set schedulerModule=swarm` to use the swarm scheduler.")
	return nil
}

//StopSwarm stops swarm
func StopSwarm(b backend.Backend) error {

	outchan := make(chan string)
	errchan := make(chan error)
	defer close(outchan)
	defer close(errchan)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Stopping Swarm...")
	outchan <- fmt.Sprintf("swarm nodes and swarm manager")
	b.Stop([]string{"swarm-node", "swarm-manager"}, &wg, outchan, errchan)
	wg.Wait()

	fmt.Println("Done.")
	fmt.Println()
	return nil
}

//UnInstallSwarm uninstall Swarm
func UnInstallSwarm(b backend.Backend) error {
	outchan := make(chan string)
	errchan := make(chan error)
	defer close(outchan)
	defer close(errchan)
	var wg sync.WaitGroup
	go printState(outchan, errchan, 500*time.Millisecond)
	outchan <- utils.DeisIfy("Destroying Swarm...")
	outchan <- fmt.Sprintf("swarm nodes and swarm manager...")
	b.Destroy([]string{"swarm-node", "swarm-manager"}, &wg, outchan, errchan)
	wg.Wait()
	fmt.Println("Done.")
	fmt.Println()
	return nil
}
