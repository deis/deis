package syslogd

import (
	"fmt"
	"log"
	"os"

	"github.com/deis/deis/logger/handler"
	"github.com/deis/deis/logger/syslog"
)

var LogRoot string
var HandlerType string
var RingBufferSize int
var WebServicePort int

// Listen starts a new syslog server which runs until it receives a signal.
func Listen(exitChan, cleanupDone chan bool, drainChan chan string, bindAddr string) {
	fmt.Println("Starting syslog...")
	// If LogRoot doesn't exist, create it
	// equivalent to Python's `if not os.path.exists(filename)`
	if _, err := os.Stat(LogRoot); os.IsNotExist(err) {
		if err = os.MkdirAll(LogRoot, 0777); err != nil {
			log.Fatalf("unable to create LogRoot at %s: %v", LogRoot, err)
		}
	}
	// Create a server with one handler and run one listen goroutine
	s := syslog.NewServer()
	var h *handler.Handler
	if HandlerType == "standard" {
		fmt.Println("Using standard handler for syslog")
		h = handler.StandardHandler(LogRoot)
	} else {
		fmt.Println("Using ring buffer handler for syslog")
		h = handler.RingBufferHandler(RingBufferSize, WebServicePort)
	}
	s.AddHandler(h)
	s.Listen(bindAddr)
	fmt.Println("Syslog server started...")
	fmt.Println("deis-logger running")

	// Wait for terminating signal
	for {
		select {
		case <-exitChan:
			// Shutdown the server
			fmt.Println("Shutting down...")
			s.Shutdown()
			cleanupDone <- true
		case d := <-drainChan:
			h.DrainURI = d
		}
	}
}
