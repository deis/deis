package handler

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

var logRootPath string

type Handler struct {
	// To simplify implementation of our handler we embed helper
	// syslog.BaseHandler struct.
	*syslog.BaseHandler
	DrainURI string
}

// Simple filter for named/bind messages which can be used with BaseHandler
func filter(m syslog.SyslogMessage) bool {
	return true
}

func StandardHandler(logRoot string) *Handler {
	logRootPath = logRoot
	h := Handler{
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

func getProjectName(message string) (string, error) {
	r := regexp.MustCompile(`^.* ([-_a-z0-9]+)\[[a-z0-9-_\.]+\].*`)
	match := r.FindStringSubmatch(message)
	if match == nil {
		return "", fmt.Errorf("Could not find app name in message: %s", message)
	}

	return match[1], nil
}

func getLogFile(message string) (io.Writer, error) {
	appName, err := getProjectName(message)
	if err != nil {
		return nil, err
	}
	filePath := path.Join(logRootPath, appName+".log")
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
func (h *Handler) mainLoop() {
	for {
		m := h.Get()
		if m == nil {
			break
		}
		if h.DrainURI != "" {
			drain.SendToDrain(m.String(), h.DrainURI)
		}
		err := writeToDisk(m)
		if err != nil {
			log.Println(err)
		}
	}
	h.End()
}
