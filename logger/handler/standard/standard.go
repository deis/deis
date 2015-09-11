package standard

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type handler struct {
	logRoot string
}

func NewHandler(logRoot string) (*handler, error) {
	src, err := os.Stat(logRoot)
	if err != nil {
		return nil, fmt.Errorf("Directory %s does exist", logRoot)
	}
	if !src.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", logRoot)
	}
	return &handler{logRoot: logRoot}, nil
}

// Write adds a log message to to an app-specific log file
func (h *handler) Write(app string, message string) error {
	f, err := h.getFile(app)
	if err != nil {
		return err
	}
	_, err = f.WriteString(message + "\n")
	if err != nil {
		return err
	}
	return nil
}

// Read retrieves a specified number of log lines from an app-specific log file
func (h *handler) Read(app string, lines int) ([]string, error) {
	if lines <= 0 {
		return []string{}, nil
	}
	filePath := h.getFilePath(app)
	// check if file exists
	exists, err := fileExists(filePath)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Could not find logs for '%s'", app)
	}
	logBytes, err := exec.Command("tail", "-n", strconv.Itoa(lines), filePath).Output()
	if err != nil {
		return nil, err
	}
	logStrs := strings.Split(string(logBytes), "\n")
	return logStrs[:len(logStrs)-1], nil
}

func (h *handler) getFile(app string) (*os.File, error) {
	filePath := h.getFilePath(app)
	// check if file exists
	exists, err := fileExists(filePath)
	if err != nil {
		return nil, err
	}
	// return a new file or the existing file for appending
	var file *os.File
	if exists {
		file, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0644)
	} else {
		file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	}
	return file, err
}

func (h *handler) getFilePath(app string) string {
	return path.Join(h.logRoot, app+".log")
}

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
