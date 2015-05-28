package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/utils"
)

// InstallMesos loads all Mesos units for StartMesos
func InstallMesos(b backend.Backend) error {

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Installing Mesos...")

	installMesosServices(b, &wg, outchan, errchan)

	wg.Wait()
	close(outchan)

	fmt.Println("Done.")
	fmt.Println()
	fmt.Println("Please run `deisctl start mesos` to boot up Mesos.")
	return nil
}

func installMesosServices(b backend.Backend, wg *sync.WaitGroup, outchan chan string, errchan chan error) {

	outchan <- fmt.Sprintf("Zookeeper...")
	b.Create([]string{"mesos-zk"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Mesos Master...")
	b.Create([]string{"mesos-master"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Mesos Slave...")
	b.Create([]string{"mesos-slave"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Marathon Framework...")
	b.Create([]string{"mesos-marathon"}, wg, outchan, errchan)
	wg.Wait()
}

// UninstallMesos unloads and uninstalls all Mesos component definitions
func UninstallMesos(b backend.Backend) error {

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Uninstalling Mesos...")

	uninstallMesosServices(b, &wg, outchan, errchan)

	wg.Wait()
	close(outchan)

	fmt.Println("Done.")
	return nil
}

func uninstallMesosServices(b backend.Backend, wg *sync.WaitGroup, outchan chan string, errchan chan error) error {

	outchan <- fmt.Sprintf("Marathon Framework...")
	b.Destroy([]string{"mesos-marathon"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Mesos Slave...")
	b.Destroy([]string{"mesos-slave"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Mesos Master...")
	b.Destroy([]string{"mesos-master"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Zookeeper...")
	b.Destroy([]string{"mesos-zk"}, wg, outchan, errchan)
	wg.Wait()

	return nil
}

// StartMesos activates all Mesos components.
func StartMesos(b backend.Backend) error {

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Starting Mesos...")

	startMesosServices(b, &wg, outchan, errchan)

	wg.Wait()
	close(outchan)

	fmt.Println("Done.")
	fmt.Println()
	fmt.Println("Please use `deisctl config controller set schedulerModule=mesos_marathon`")
	return nil
}

func startMesosServices(b backend.Backend, wg *sync.WaitGroup, outchan chan string, errchan chan error) {

	outchan <- fmt.Sprintf("Zookeeper...")
	b.Start([]string{"mesos-zk"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Mesos Master...")
	b.Start([]string{"mesos-master"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Mesos Slave...")
	b.Start([]string{"mesos-slave"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Marathon Framework...")
	b.Start([]string{"mesos-marathon"}, wg, outchan, errchan)
	wg.Wait()
}

// StopMesos deactivates all Mesos components.
func StopMesos(b backend.Backend) error {

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Stopping Mesos...")

	stopMesosServices(b, &wg, outchan, errchan)

	wg.Wait()
	close(outchan)

	fmt.Println("Done.")
	fmt.Println()
	fmt.Println("Please run `deisctl start mesos` to restart Mesos.")
	return nil
}

func stopMesosServices(b backend.Backend, wg *sync.WaitGroup, outchan chan string, errchan chan error) {

	outchan <- fmt.Sprintf("Marathon Framework...")
	b.Stop([]string{"mesos-marathon"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Mesos Slave...")
	b.Stop([]string{"mesos-slave"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Mesos Master...")
	b.Stop([]string{"mesos-master"}, wg, outchan, errchan)
	wg.Wait()

	outchan <- fmt.Sprintf("Zookeeper...")
	b.Stop([]string{"mesos-zk"}, wg, outchan, errchan)
	wg.Wait()
}
