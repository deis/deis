// Package basher provides an API for running and integrating with Bash from Go
package basher

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/kardianos/osext"
)

func exitStatus(err error) (int, error) {
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// There is no platform independent way to retrieve
			// the exit code, but the following will work on Unix
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return int(status.ExitStatus()), nil
			}
		}
		return 0, err
	}
	return 0, nil
}

// Application sets up a common entrypoint for a Bash application that
// uses exported Go functions. It uses the DEBUG environment variable
// to set debug on the Context, and SHELL for the Bash binary if it
// includes the string "bash". You can pass a loader function to use
// for the sourced files, and a boolean for whether or not the
// environment should be copied into the Context process.
func Application(
	funcs map[string]func([]string),
	scripts []string,
	loader func(string) ([]byte, error),
	copyEnv bool) {

	var bashPath string
	bashPath, err := exec.LookPath("bash")
	if err != nil {
		if strings.Contains(os.Getenv("SHELL"), "bash") {
			bashPath = os.Getenv("SHELL")
		} else {
			bashPath = "/bin/bash"
		}
	}
	bash, err := NewContext(bashPath, os.Getenv("DEBUG") != "")
	if err != nil {
		log.Fatal(err)
	}
	for name, fn := range funcs {
		bash.ExportFunc(name, fn)
	}
	if bash.HandleFuncs(os.Args) {
		os.Exit(0)
	}

	for _, script := range scripts {
		bash.Source(script, loader)
	}
	if copyEnv {
		bash.CopyEnv()
	}
	status, err := bash.Run("main", os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(status)
}

// A Context is an instance of a Bash interpreter and environment, including
// sourced scripts, environment variables, and embedded Go functions
type Context struct {
	sync.Mutex

	// Debug simply leaves the generated BASH_ENV file produced
	// from each Run call of this Context for debugging.
	Debug bool

	// BashPath is the path to the Bash executable to be used by Run
	BashPath string

	// SelfPath is set by NewContext to be the current executable path.
	// It's used to call back into the calling Go process to run exported
	// functions.
	SelfPath string

	// The io.Reader given to Bash for STDIN
	Stdin io.Reader

	// The io.Writer given to Bash for STDOUT
	Stdout io.Writer

	// The io.Writer given to Bash for STDERR
	Stderr io.Writer

	vars    []string
	scripts [][]byte
	funcs   map[string]func([]string)
}

// Creates and initializes a new Context that will use the given Bash executable.
// The debug mode will leave the produced temporary BASH_ENV file for inspection.
func NewContext(bashpath string, debug bool) (*Context, error) {
	executable, err := osext.Executable()
	if err != nil {
		return nil, err
	}
	return &Context{
		Debug:    debug,
		BashPath: bashpath,
		SelfPath: executable,
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		scripts:  make([][]byte, 0),
		vars:     make([]string, 0),
		funcs:    make(map[string]func([]string)),
	}, nil
}

// Copies the current environment variables into the Context
func (c *Context) CopyEnv() {
	c.Lock()
	defer c.Unlock()
	c.vars = append(c.vars, os.Environ()...)
}

// Adds a shell script to the Context environment. The loader argument can be nil
// which means it will use ioutil.Readfile and load from disk, but it exists so you
// can use the Asset function produced by go-bindata when including script files in
// your Go binary. Calls to Source adds files to the environment in order.
func (c *Context) Source(filepath string, loader func(string) ([]byte, error)) error {
	if loader == nil {
		loader = ioutil.ReadFile
	}
	data, err := loader(filepath)
	if err != nil {
		return err
	}
	c.Lock()
	defer c.Unlock()
	c.scripts = append(c.scripts, data)
	return nil
}

// Adds an environment variable to the Context
func (c *Context) Export(name string, value string) {
	c.Lock()
	defer c.Unlock()
	c.vars = append(c.vars, name+"="+value)
}

// Registers a function with the Context that will produce a Bash function in the environment
// that calls back into your executable triggering the function defined as fn.
func (c *Context) ExportFunc(name string, fn func([]string)) {
	c.Lock()
	defer c.Unlock()
	c.funcs[name] = fn
}

// Expects your os.Args to parse and handle any callbacks to Go functions registered with
// ExportFunc. You normally call this at the beginning of your program. If a registered
// function is found and handled, HandleFuncs will exit with the appropriate exit code for you.
func (c *Context) HandleFuncs(args []string) bool {
	for i, arg := range args {
		if arg == "::" && len(args) > i+1 {
			c.Lock()
			defer c.Unlock()
			for cmd := range c.funcs {
				if cmd == args[i+1] {
					c.funcs[cmd](args[i+2:])
					return true
				}
			}
			return false
		}
	}
	return false
}

func (c *Context) buildEnvfile() (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "bashenv.")
	if err != nil {
		return "", err
	}
	defer file.Close()
	// variables
	file.Write([]byte("unset BASH_ENV\n")) // unset for future calls to bash
	file.Write([]byte("export SELF=" + os.Args[0] + "\n"))
	file.Write([]byte("export EXECUTABLE='" + c.SelfPath + "'\n"))
	for _, kvp := range c.vars {
		file.Write([]byte("export " + strings.Replace(
			strings.Replace(kvp, "'", "\\'", -1), "=", "=$'", 1) + "'\n"))
	}
	// functions
	for cmd := range c.funcs {
		file.Write([]byte(cmd + "() { $EXECUTABLE :: " + cmd + " \"$@\"; }\n"))
	}
	// scripts
	for _, data := range c.scripts {
		file.Write(append(data, '\n'))
	}
	return file.Name(), nil
}

// Runs a command in Bash from this Context. With each call, a temporary file
// is generated used as BASH_ENV when calling Bash that includes all variables,
// sourced scripts, and exported functions from the Context. Standard I/O by
// default is attached to the calling process I/O. You can change this by setting
// the Stdout, Stderr, Stdin variables of the Context.
func (c *Context) Run(command string, args []string) (int, error) {
	c.Lock()
	defer c.Unlock()
	envfile, err := c.buildEnvfile()
	if err != nil {
		return 0, err
	}
	if !c.Debug {
		defer os.Remove(envfile)
	}
	argstring := ""
	for _, arg := range args {
		argstring = argstring + " '" + arg + "'"
	}
	cmd := exec.Command(c.BashPath, "-c", command+argstring)
	cmd.Env = []string{"BASH_ENV=" + envfile}
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	return exitStatus(cmd.Run())
}
