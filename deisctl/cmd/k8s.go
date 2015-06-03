package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/utils"
)

//InstallK8s Installs K8s
func InstallK8s(b backend.Backend) error {
	outchan := make(chan string)
	errchan := make(chan error)
	defer close(outchan)
	defer close(errchan)
	var wg sync.WaitGroup
	go printState(outchan, errchan, 500*time.Millisecond)
	outchan <- utils.DeisIfy("Installing K8s...")
	outchan <- fmt.Sprintf("K8s API Server ...")
	b.Create([]string{"kube-apiserver"}, &wg, outchan, errchan)
	wg.Wait()
  outchan <- fmt.Sprintf("K8s controller and scheduler ...")
  b.Create([]string{"kube-controller-manager","kube-scheduler"}, &wg, outchan, errchan)
  wg.Wait()
  outchan <- fmt.Sprintf("K8s proxy and kubelet ...")
  b.Create([]string{"kube-proxy","kube-kubelet"}, &wg, outchan, errchan)
  wg.Wait()
  fmt.Println("Done.")
	fmt.Println()
	fmt.Println("Please run `deisctl start k8s` to start K8s.")
	return nil
}

//StartK8s starts K8s Schduler
func StartK8s(b backend.Backend) error {
	outchan := make(chan string)
	errchan := make(chan error)
	defer close(outchan)
	defer close(errchan)
	var wg sync.WaitGroup
	go printState(outchan, errchan, 500*time.Millisecond)
	outchan <- utils.DeisIfy("Starting K8s...")
  outchan <- fmt.Sprintf("K8s API Server ...")
	b.Start([]string{"kube-apiserver"}, &wg, outchan, errchan)
	wg.Wait()
  outchan <- fmt.Sprintf("K8s controller and scheduler ...")
  b.Start([]string{"kube-controller-manager","kube-scheduler"}, &wg, outchan, errchan)
  wg.Wait()
  outchan <- fmt.Sprintf("K8s proxy and kubelet ...")
  b.Start([]string{"kube-proxy","kube-kubelet"}, &wg, outchan, errchan)
  wg.Wait()
	fmt.Println("Done.")
	fmt.Println("Please run `deisctl config controller set schedulerModule=k8s` to use the K8s scheduler.")
	return nil
}

//StopK8s stops K8s
func StopK8s(b backend.Backend) error {

	outchan := make(chan string)
	errchan := make(chan error)
	defer close(outchan)
	defer close(errchan)
	var wg sync.WaitGroup

	go printState(outchan, errchan, 500*time.Millisecond)

	outchan <- utils.DeisIfy("Stopping K8s...")
  outchan <- fmt.Sprintf("K8s proxy and kubelet ...")
  b.Stop([]string{"kube-proxy","kube-kubelet"}, &wg, outchan, errchan)
  wg.Wait()
  outchan <- fmt.Sprintf("K8s controller and scheduler ...")
  b.Stop([]string{"kube-controller-manager","kube-scheduler"}, &wg, outchan, errchan)
  wg.Wait()
  outchan <- fmt.Sprintf("K8s API Server ...")
	b.Stop([]string{"kube-apiserver"}, &wg, outchan, errchan)
	wg.Wait()
	fmt.Println("Done.")
	fmt.Println()
	return nil
}

//UnInstallK8s uninstall K8s
func UnInstallK8s(b backend.Backend) error {
	outchan := make(chan string)
	errchan := make(chan error)
	defer close(outchan)
	defer close(errchan)
	var wg sync.WaitGroup
	go printState(outchan, errchan, 500*time.Millisecond)
	outchan <- utils.DeisIfy("Destroying K8s...")
  outchan <- fmt.Sprintf("K8s proxy and kubelet ...")
  b.Destroy([]string{"kube-proxy","kube-kubelet"}, &wg, outchan, errchan)
  wg.Wait()
  outchan <- fmt.Sprintf("K8s controller and scheduler ...")
  b.Destroy([]string{"kube-controller-manager","kube-scheduler"}, &wg, outchan, errchan)
  wg.Wait()
  outchan <- fmt.Sprintf("K8s API Server ...")
	b.Destroy([]string{"kube-apiserver"}, &wg, outchan, errchan)
	wg.Wait()
	fmt.Println("Done.")
	fmt.Println()
	return nil
}
