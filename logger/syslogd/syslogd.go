package syslogd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"

	"github.com/deis/deis/logger/syslog"
)

const logRoot = "/data/logs"

type handler struct {
	// To simplify implementation of our handler we embed helper
	// syslog.BaseHandler struct.
	*syslog.BaseHandler
}

// Simple fiter for named/bind messages which can be used with BaseHandler
func filter(m syslog.SyslogMessage) bool {
	return true
}

func newHandler() *handler {
	h := handler{syslog.NewBaseHandler(5, filter, false)}
	go h.mainLoop() // BaseHandler needs some gorutine that reads from its queue
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
	filePath := path.Join(logRoot, appName+".log")
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
		err := writeToDisk(m)
		if err != nil {
			log.Println(err)
		}
	}
	h.End()
}

// Listen starts a new syslog server which runs until it receives a signal.
func Listen(signalChan chan os.Signal, cleanupDone chan bool) {
	fmt.Println("Starting syslog...")
	// Create a server with one handler and run one listen gorutine
	s := syslog.NewServer()
	s.AddHandler(newHandler())
	s.Listen("0.0.0.0:514")
	fmt.Println("Syslog server started...")
	fmt.Println("deis-logger running")

	// Wait for terminating signal
	for _ = range signalChan {
		// Shutdown the server
		fmt.Println("Shutting down...")
		s.Shutdown()
		cleanupDone <- true
	}
}
