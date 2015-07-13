package confd

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

const (
	confdError = "Apr 29 17:00:54 deis-1 sh: 2015-04-29T15:00:54Z fceb2accfbf5 confd[1484]: ERROR template: builder:59:47: executing \"builder\" at <getv \"/deis/registry...>: error calling getv: key does not exist"
)

func TestReturnError(t *testing.T) {
	signalChan := make(chan os.Signal, 1)
	args := "while(true);do echo '" + confdError + "'; sleep 1;done"
	cmd := exec.Command("/bin/bash", "-c", args)

	stdout, err := cmd.StdoutPipe()
	checkError(signalChan, err)

	cmd.Start()

	go checkNumberOfErrors(stdout, 1, 2*time.Second, signalChan)

	for {
		select {
		case <-time.Tick(5 * time.Second):
			return
		case s := <-signalChan:
			log.Debugf("Signal received: %v", s)
			switch s {
			case syscall.SIGKILL:
				// we expect this
				return
			}
		}
	}
}

func TestReturnWithoutError(t *testing.T) {
	signalChan := make(chan os.Signal, 1)
	args := "while(true);do echo '" + confdError + "'; sleep 1;done"
	cmd := exec.Command("/bin/bash", "-c", args)

	stdout, err := cmd.StdoutPipe()
	checkError(signalChan, err)

	cmd.Start()

	go checkNumberOfErrors(stdout, 2, 2*time.Second, signalChan)

	for {
		select {
		case <-time.Tick(5 * time.Second):
			return
		case s := <-signalChan:
			log.Debugf("Signal received: %v", s)
			switch s {
			case syscall.SIGKILL:
				t.Fatal("Unexpected error received")
			}
		}
	}
}
