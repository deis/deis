package syslogd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"

	"github.com/deis/deis/logger/syslog"

	"github.com/deis/deis/logger/drain"
)

// LogRoot is the log path to store logs.
var LogRoot string

type handler struct {
	// To simplify implementation of our handler we embed helper
	// syslog.BaseHandler struct.
	*syslog.BaseHandler
	drainURI string
}

// Simple fiter for named/bind messages which can be used with BaseHandler
func filter(m syslog.SyslogMessage) bool {
	return true
}

func newHandler() *handler {
	h := handler{
		BaseHandler: syslog.NewBaseHandler(5, filter, false),
	}

	go h.mainLoop() // BaseHandler needs some goroutine that reads from its queue
	return &h
}

// check if a file path exists
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getLogFile(message string) (io.Writer, error) {
	r := regexp.MustCompile(`^.* ([-a-z0-9]+)\[[a-z0-9-_\.]+\].*`)
	match := r.FindStringSubmatch(message)
	if match == nil {
		return nil, fmt.Errorf("Could not find app name in message: %s", message)
	}
	appName := match[1]
	filePath := path.Join(LogRoot, appName+".log")
	// check if file exists
	exists, err := fileExists(filePath)
	if err != nil {
		return nil, err
	}
	// return a new file or the existing file for appending
	var file io.Writer
	if exists {
		file, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0644)
	} else {
		file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	}
	return file, err
}

func writeToDisk(m syslog.SyslogMessage) error {
	file, err := getLogFile(m.String())
	if err != nil {
		return err
	}
	bytes := []byte(m.String() + "\n")
	file.Write(bytes)
	return nil
}

// mainLoop reads from BaseHandler queue using h.Get and logs messages to stdout
func (h *handler) mainLoop() {
	for {
		m := h.Get()
		if m == nil {
			break
		}
		if h.drainURI != "" {
			drain.SendToDrain(m.String(), h.drainURI)
		}
		err := writeToDisk(m)
		if err != nil {
			log.Println(err)
		}
	}
	h.End()
}

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
	h := newHandler()
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
			h.drainURI = d
		}
	}
}
