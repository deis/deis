package file

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
)

type adapter struct {
	logRoot string
	files   map[string]*os.File
	mutex   sync.Mutex
}

// NewStorageAdapter returns a pointer to a new instance of a file-based storage.Adapter.
func NewStorageAdapter(logRoot string) (*adapter, error) {
	src, err := os.Stat(logRoot)
	if err != nil {
		return nil, fmt.Errorf("Directory %s does not exist", logRoot)
	}
	if !src.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", logRoot)
	}
	return &adapter{logRoot: logRoot, files: make(map[string]*os.File)}, nil
}

// Write adds a log message to to an app-specific log file
func (a *adapter) Write(app string, message string) error {
	// Check first if we might actually have to add to the map of file pointers so we can avoid
	// waiting for / obtaining a lock unnecessarily
	f, ok := a.files[app]
	if !ok {
		// Ensure only one goroutine at a time can be adding a file pointer to the map of file
		// pointers
		a.mutex.Lock()
		defer a.mutex.Unlock()
		f, ok = a.files[app]
		if !ok {
			var err error
			f, err = a.getFile(app)
			if err != nil {
				return err
			}
			a.files[app] = f
		}
	}
	if _, err := f.WriteString(message + "\n"); err != nil {
		return err
	}
	return nil
}

// Read retrieves a specified number of log lines from an app-specific log file
func (a *adapter) Read(app string, lines int) ([]string, error) {
	if lines <= 0 {
		return []string{}, nil
	}
	filePath := a.getFilePath(app)
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

// Destroy deletes stored logs for the specified application
func (a *adapter) Destroy(app string) error {
	// Check first if the map of file pointers even contains the file pointer we want so we can avoid
	// waiting for / obtaining a lock unnecessarily
	f, ok := a.files[app]
	if ok {
		// Ensure no other goroutine is trying to modify the file pointer map while we're trying to
		// clean up
		a.mutex.Lock()
		defer a.mutex.Unlock()
		exists, err := fileExists(f.Name())
		if err != nil {
			return err
		}
		if exists {
			if err := os.Remove(f.Name()); err != nil {
				return err
			}
		}
		delete(a.files, app)
	}
	return nil
}

func (a *adapter) Reopen() error {
	// Ensure no other goroutine is trying to add a file pointer to the map of file pointers while
	// we're trying to clear it out
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.files = make(map[string]*os.File)
	return nil
}

func (a *adapter) getFile(app string) (*os.File, error) {
	filePath := a.getFilePath(app)
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

func (a *adapter) getFilePath(app string) string {
	return path.Join(a.logRoot, app+".log")
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
