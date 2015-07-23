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
	fmt.Fprintln(Stdout, "K8s API Server...")
	b.Create([]string{"kube-apiserver"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s controller and scheduler ...")
	b.Create([]string{"kube-controller-manager", "kube-scheduler"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s proxy and kubelet ...")
	b.Create([]string{"kube-proxy", "kube-kubelet"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "Done.\n ")
	fmt.Fprintln(Stdout, "Please run `deisctl start k8s` to start K8s.")
	return nil
}

//StartK8s starts K8s Schduler
func StartK8s(b backend.Backend) error {
	var wg sync.WaitGroup
	io.WriteString(Stdout, prettyprint.DeisIfy("Starting K8s..."))
	fmt.Fprintln(Stdout, "K8s API Server ...")
	b.Start([]string{"kube-apiserver"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s controller and scheduler ...")
	b.Start([]string{"kube-controller-manager", "kube-scheduler"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s proxy and kubelet ...")
	b.Start([]string{"kube-proxy", "kube-kubelet"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "Done.\n ")
	fmt.Fprintln(Stdout, "Please run `deisctl config controller set schedulerModule=k8s` to use the K8s scheduler.")
	return nil
}

//StopK8s stops K8s
func StopK8s(b backend.Backend) error {
	var wg sync.WaitGroup
	io.WriteString(Stdout, prettyprint.DeisIfy("Stopping K8s..."))
	fmt.Fprintln(Stdout, "K8s proxy and kubelet ...")
	b.Stop([]string{"kube-proxy", "kube-kubelet"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s controller and scheduler ...")
	b.Stop([]string{"kube-controller-manager", "kube-scheduler"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s API Server ...")
	b.Stop([]string{"kube-apiserver"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "Done.\n ")
	return nil
}

//UnInstallK8s uninstall K8s
func UnInstallK8s(b backend.Backend) error {
	var wg sync.WaitGroup
	io.WriteString(Stdout, prettyprint.DeisIfy("Destroying K8s..."))
	fmt.Fprintln(Stdout, "K8s proxy and kubelet ...")
	b.Destroy([]string{"kube-proxy", "kube-kubelet"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s controller and scheduler ...")
	b.Destroy([]string{"kube-controller-manager", "kube-scheduler"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "K8s API Server ...")
	b.Destroy([]string{"kube-apiserver"}, &wg, Stdout, Stderr)
	wg.Wait()
	fmt.Fprintln(Stdout, "Done.\n ")
	return nil
}
