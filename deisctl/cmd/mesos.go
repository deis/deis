package cmd

import (
	"fmt"
	"io"
	"sync"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/pkg/prettyprint"
)

// InstallMesos loads all Mesos units for StartMesos
func InstallMesos(b backend.Backend) error {

	var wg sync.WaitGroup

	io.WriteString(Stdout, prettyprint.DeisIfy("Installing Mesos..."))

	installMesosServices(b, &wg, Stdout, Stderr)

	wg.Wait()

	fmt.Fprintln(Stdout, "Done.")
	fmt.Fprintln(Stdout, "")
	fmt.Fprintln(Stdout, "Please run `deisctl start mesos` to boot up Mesos.")
	return nil
}

func installMesosServices(b backend.Backend, wg *sync.WaitGroup, out, err io.Writer) {

	fmt.Fprintln(out, "Zookeeper...")
	b.Create([]string{"zookeeper"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Mesos Master...")
	b.Create([]string{"mesos-master"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Mesos Slave...")
	b.Create([]string{"mesos-slave"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Marathon Framework...")
	b.Create([]string{"mesos-marathon"}, wg, out, err)
	wg.Wait()
}

// UninstallMesos unloads and uninstalls all Mesos component definitions
func UninstallMesos(b backend.Backend) error {

	var wg sync.WaitGroup

	io.WriteString(Stdout, prettyprint.DeisIfy("Uninstalling Mesos..."))

	uninstallMesosServices(b, &wg, Stdout, Stderr)

	wg.Wait()

	fmt.Fprintln(Stdout, "Done.")
	return nil
}

func uninstallMesosServices(b backend.Backend, wg *sync.WaitGroup, out, err io.Writer) error {

	fmt.Fprintln(out, "Marathon Framework...")
	b.Destroy([]string{"mesos-marathon"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Mesos Slave...")
	b.Destroy([]string{"mesos-slave"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Mesos Master...")
	b.Destroy([]string{"mesos-master"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Zookeeper...")
	b.Destroy([]string{"zookeeper"}, wg, out, err)
	wg.Wait()

	return nil
}

// StartMesos activates all Mesos components.
func StartMesos(b backend.Backend) error {

	var wg sync.WaitGroup

	io.WriteString(Stdout, prettyprint.DeisIfy("Starting Mesos..."))

	startMesosServices(b, &wg, Stdout, Stderr)

	wg.Wait()

	fmt.Fprintln(Stdout, "Done.")
	fmt.Fprintln(Stdout, "")
	fmt.Fprintln(Stdout, "Please use `deisctl config controller set schedulerModule=mesos_marathon`")
	return nil
}

func startMesosServices(b backend.Backend, wg *sync.WaitGroup, out, err io.Writer) {

	fmt.Fprintln(out, "Zookeeper...")
	b.Start([]string{"zookeeper"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Mesos Master...")
	b.Start([]string{"mesos-master"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Mesos Slave...")
	b.Start([]string{"mesos-slave"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Marathon Framework...")
	b.Start([]string{"mesos-marathon"}, wg, out, err)
	wg.Wait()
}

// StopMesos deactivates all Mesos components.
func StopMesos(b backend.Backend) error {

	var wg sync.WaitGroup

	io.WriteString(Stdout, prettyprint.DeisIfy("Stopping Mesos..."))

	stopMesosServices(b, &wg, Stdout, Stderr)

	wg.Wait()

	fmt.Fprintln(Stdout, "Done.")
	fmt.Fprintln(Stdout, "")
	fmt.Fprintln(Stdout, "Please run `deisctl start mesos` to restart Mesos.")
	return nil
}

func stopMesosServices(b backend.Backend, wg *sync.WaitGroup, out, err io.Writer) {

	fmt.Fprintln(out, "Marathon Framwork...")
	b.Stop([]string{"mesos-marathon"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Mesos Slave...")
	b.Stop([]string{"mesos-slave"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Mesos Master...")
	b.Stop([]string{"mesos-master"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Zookeeper...")
	b.Stop([]string{"zookeeper"}, wg, out, err)
	wg.Wait()
}
