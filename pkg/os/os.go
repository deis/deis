package os

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	Log "github.com/deis/deis/pkg/log"

	"github.com/progrium/go-basher"
)

var log = Log.New()

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
		log.Printf("an error ocurred executing command: [%s params %v], %v", command, args, err)
		signalChan <- syscall.SIGTERM
	}

	err = cmd.Wait()
	log.Printf("command finished with error: %v", err)
	signalChan <- syscall.SIGTERM
}

// RunScript run a shell script using go-basher and if it returns an error
// send a signal to terminate the execution
func RunScript(signalChan chan os.Signal, script string, params map[string]string,
	loader func(string) ([]byte, error)) {
	bash, _ := basher.NewContext("/bin/bash", false)
	bash.Source(script, loader)
	if params != nil {
		for key, value := range params {
			bash.Export(key, value)
		}
	}

	_, err := bash.Run("main", nil)
	if err != nil {
		log.Printf("command finished with error: %v", err)
		signalChan <- syscall.SIGTERM
	}
}

// RunCommand run a command and return.
func RunCommand(signalChan chan os.Signal, command string, args []string, signalErrors bool) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if signalErrors {
			signalChan <- syscall.SIGTERM
		}

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
