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

	io.WriteString(Stdout, prettyprint.DeisIfy("Installing Mesos/Marathon..."))

	installMesosServices(b, &wg, Stdout, Stderr)

	wg.Wait()

	fmt.Fprintln(Stdout, "Done.\n ")
	fmt.Fprintln(Stdout, "Please run `deisctl start mesos` to boot up Mesos.")
	return nil
}

func installMesosServices(b backend.Backend, wg *sync.WaitGroup, out, err io.Writer) {

	fmt.Fprintln(out, "Mesos/Marathon control plane...")
	b.Create([]string{"zookeeper", "mesos-master"}, wg, out, err)
	wg.Wait()
	b.Create([]string{"mesos-marathon"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Mesos/Marathon data plane...")
	b.Create([]string{"mesos-slave"}, wg, out, err)
	wg.Wait()
}

// UninstallMesos unloads and uninstalls all Mesos component definitions
func UninstallMesos(b backend.Backend) error {

	var wg sync.WaitGroup

	io.WriteString(Stdout, prettyprint.DeisIfy("Uninstalling Mesos/Marathon..."))

	uninstallMesosServices(b, &wg, Stdout, Stderr)

	wg.Wait()

	fmt.Fprintln(Stdout, "Done.\n ")
	return nil
}

func uninstallMesosServices(b backend.Backend, wg *sync.WaitGroup, out, err io.Writer) error {

	fmt.Fprintln(out, "Mesos/Marathon data plane...")
	b.Destroy([]string{"mesos-slave"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Mesos/Marathon control plane...")
	b.Destroy([]string{"mesos-marathon", "mesos-master", "zookeeper"}, wg, out, err)
	wg.Wait()

	return nil
}

// StartMesos activates all Mesos components.
func StartMesos(b backend.Backend, verbose bool) error {

	var wg sync.WaitGroup

	io.WriteString(Stdout, prettyprint.DeisIfy("Starting Mesos/Marathon..."))

	startMesosServices(b, verbose, &wg, Stdout, Stderr)

	wg.Wait()

	fmt.Fprintln(Stdout, "Done.\n ")
	fmt.Fprintln(Stdout, "Please use `deisctl config controller set schedulerModule=mesos_marathon`")
	return nil
}

func startMesosServices(b backend.Backend, verbose bool, wg *sync.WaitGroup, out, err io.Writer) {
	var quitChans []chan bool

	fmt.Fprintln(out, "Mesos/Marathon control plane...")
	b.Start([]string{"zookeeper"}, wg, out, err)
	wg.Wait()
	b.Start([]string{"mesos-master"}, wg, out, err)
	startLogging(b, verbose, &quitChans, "mesos-master", Stdout)
	wg.Wait()
	stopLogging(&quitChans)
	b.Start([]string{"mesos-marathon"}, wg, out, err)
	startLogging(b, verbose, &quitChans, "meos-marathon", Stdout)
	wg.Wait()
	stopLogging(&quitChans)
	fmt.Fprintln(out, "Mesos/Marathon data plane...")
	b.Start([]string{"mesos-slave"}, wg, out, err)
	wg.Wait()
}

// StopMesos deactivates all Mesos components.
func StopMesos(b backend.Backend) error {

	var wg sync.WaitGroup

	io.WriteString(Stdout, prettyprint.DeisIfy("Stopping Mesos/Marathon..."))

	stopMesosServices(b, &wg, Stdout, Stderr)

	wg.Wait()

	fmt.Fprintln(Stdout, "Done.\n ")
	fmt.Fprintln(Stdout, "Please run `deisctl start mesos` to restart Mesos.")
	return nil
}

func stopMesosServices(b backend.Backend, wg *sync.WaitGroup, out, err io.Writer) {

	fmt.Fprintln(out, "Mesos/Marathon data plane...")
	b.Stop([]string{"mesos-slave"}, wg, out, err)
	wg.Wait()

	fmt.Fprintln(out, "Mesos/Marathon control plane...")
	b.Stop([]string{"mesos-marathon"}, wg, out, err)
	wg.Wait()
	b.Stop([]string{"mesos-master"}, wg, out, err)
	wg.Wait()
	b.Stop([]string{"zookeeper"}, wg, out, err)
	wg.Wait()
}
