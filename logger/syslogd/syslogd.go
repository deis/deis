package main

import (
	"errors"
	"fmt"
	"github.com/deis/deis/logger/syslog"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"regexp"
	"syscall"
)

const logRoot = "/var/log/deis"

type handler struct {
	// To simplify implementation of our handler we embed helper
	// syslog.BaseHandler struct.
	*syslog.BaseHandler
}

// Simple fiter for named/bind messages which can be used with BaseHandler
func filter(m *syslog.Message) bool {
	// return m.Tag == "named" || m.Tag == "bind"
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

func getLogFile(m *syslog.Message) (io.Writer, error) {
	r := regexp.MustCompile(`^.* ([-a-z0-9]+)\[[a-z0-9\.]+\].*`)
	match := r.FindStringSubmatch(m.String())
	if match == nil {
		return nil, errors.New("Could not find app name in message")
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

func writeToDisk(m *syslog.Message) error {
	file, err := getLogFile(m)
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
		fmt.Println(m)
		err := writeToDisk(m)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Exit handler")
	h.End()
}

func webHandler(response http.ResponseWriter, request *http.Request) {
	appName := request.URL.Path[1:]
	r := regexp.MustCompile(`^[-a-z0-9]+`)
	match := r.FindStringSubmatch(appName)
	if match == nil {
		fmt.Fprint(response, "Could not find app name in request")
		return
	}
	filePath := path.Join(logRoot, appName+".log")
	// check if file exists
	exists, err := fileExists(filePath)
	if err != nil {
		fmt.Fprint(response, err)
		return
	}
	if !exists {
		fmt.Fprint(response, "File not found")
		return
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Fprint(response, "Could not read from file")
		return
	}
	fmt.Fprint(response, string(data))
}

func main() {
	// Create a server with one handler and run one listen gorutine
	s := syslog.NewServer()
	s.AddHandler(newHandler())
	s.Listen("0.0.0.0:514")

	http.HandleFunc("/", webHandler)
	http.ListenAndServe(":1337", nil)

	// Wait for terminating signal
	sc := make(chan os.Signal, 2)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT)
	<-sc

	// Shutdown the server
	fmt.Println("Shutdown the server...")
	s.Shutdown()
	fmt.Println("Server is down")
}
