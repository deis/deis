package os

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	logger "github.com/deis/deis/mesos/pkg/log"
	basher "github.com/progrium/go-basher"
)

const (
	alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

var log = logger.New()

// Getopt return the value of and environment variable or a default
func Getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Debugf("returning default value \"%s\" for key \"%s\"", dfault, name)
		value = dfault
	}
	return value
}

// RunProcessAsDaemon start a child process that will run indefinitely
func RunProcessAsDaemon(signalChan chan os.Signal, command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Errorf("an error ocurred executing command: [%s params %v], %v", command, args, err)
		signalChan <- syscall.SIGKILL
	}

	err = cmd.Wait()
	log.Errorf("command finished with error: %v", err)
	signalChan <- syscall.SIGKILL
}

// RunScript run a shell script using go-basher and if it returns an error
// send a signal to terminate the execution
func RunScript(script string, params map[string]string, loader func(string) ([]byte, error)) error {
	log.Debugf("running script %v", script)
	bash, _ := basher.NewContext("/bin/bash", log.Level.String() == "debug")
	bash.Source(script, loader)
	if params != nil {
		for key, value := range params {
			bash.Export(key, value)
		}
	}

	status, err := bash.Run("main", []string{})
	if err != nil {
		return err
	}
	if status != 0 {
		return fmt.Errorf("invalid exit code running script [%v]", status)
	}

	return nil
}

// RunCommand run a command and return.
func RunCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// BuildCommandFromString parses a string containing a command and multiple
// arguments and returns a valid tuple to pass to exec.Command
func BuildCommandFromString(input string) (string, []string) {
	command := strings.Split(input, " ")

	if len(command) > 1 {
		return command[0], command[1:]
	}

	return command[0], []string{}
}

// Random return a random string
func Random(size int) (string, error) {
	if size <= 0 {
		return "", errors.New("invalid size. It must be bigger or equal to 1")
	}

	var bytes = make([]byte, size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes), nil
}

// Hostname returns the host name reported by the kernel.
func Hostname() (name string, err error) {
	return os.Hostname()
}

// CopyFile copies a source file to a destination.
func CopyFile(src string, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
