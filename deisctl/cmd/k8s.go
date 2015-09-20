package cmd

import (
	"fmt"
	"io"
	"sync"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/pkg/prettyprint"
)

//InstallK8s Installs K8s
func InstallK8s(b backend.Backend) error {
	var wg sync.WaitGroup
	io.WriteString(Stdout, prettyprint.DeisIfy("Installing K8s..."))
	fmt.Fprintln(Stdout, "K8s control plane...")
	b.Create([]string{"kube-apiserver"}, &wg, Stdout, Stderr)
	wg.Wait()
	b.Create([]string{"kube-controller-manager", "kube-scheduler"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s data plane...")
	b.Create([]string{"kube-kubelet"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s router mesh...")
	b.Create([]string{"kube-proxy"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "Done.\n ")
	fmt.Fprintln(Stdout, "Please run `deisctl start k8s` to start K8s.")
	return nil
}

//StartK8s starts K8s Schduler
func StartK8s(b backend.Backend, verbose bool) error {
	var quitChans []chan bool
	var wg sync.WaitGroup
	io.WriteString(Stdout, prettyprint.DeisIfy("Starting K8s..."))
	fmt.Fprintln(Stdout, "K8s control plane...")
	b.Start([]string{"kube-apiserver"}, &wg, Stdout, Stderr)
	startLogging(b, verbose, &quitChans, "kube-apiserver", Stdout)
	wg.Wait()
	stopLogging(&quitChans)
	b.Start([]string{"kube-controller-manager", "kube-scheduler"}, &wg, Stdout, Stderr)
	startLogging(b, verbose, &quitChans, "kube-controller-manager", Stdout)
	startLogging(b, verbose, &quitChans, "kube-scheduler", Stdout)
	wg.Wait()
	stopLogging(&quitChans)
	fmt.Fprintln(Stdout, "K8s data plane...")
	b.Start([]string{"kube-kubelet"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s router mesh...")
	b.Start([]string{"kube-proxy"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "Done.\n ")
	fmt.Fprintln(Stdout, "Please run `deisctl config controller set schedulerModule=k8s` to use the K8s scheduler.")
	return nil
}

//StopK8s stops K8s
func StopK8s(b backend.Backend) error {
	var wg sync.WaitGroup
	io.WriteString(Stdout, prettyprint.DeisIfy("Stopping K8s..."))
	fmt.Fprintln(Stdout, "K8s router mesh...")
	b.Stop([]string{"kube-proxy"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s data plane...")
	b.Stop([]string{"kube-kubelet"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s control plane...")
	b.Stop([]string{"kube-controller-manager", "kube-scheduler"}, &wg, Stdout, Stderr)
	wg.Wait()
	b.Stop([]string{"kube-apiserver"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "Done.\n ")
	return nil
}

//UnInstallK8s uninstall K8s
func UnInstallK8s(b backend.Backend) error {
	var wg sync.WaitGroup
	io.WriteString(Stdout, prettyprint.DeisIfy("Uninstalling K8s..."))
	fmt.Fprintln(Stdout, "K8s router mesh...")
	b.Destroy([]string{"kube-proxy"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s data plane...")
	b.Destroy([]string{"kube-kubelet"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s control plane...")
	b.Destroy([]string{"kube-controller-manager", "kube-scheduler"}, &wg, Stdout, Stderr)
	wg.Wait()
	b.Destroy([]string{"kube-apiserver"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "Done.\n ")
	return nil
}
