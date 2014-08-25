// Package utils contains commonly useful functions from Deis testing.

package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// Append grows a string array by appending a new element.
func Append(slice []string, data string) []string {
	m := len(slice)
	n := m + 1
	if n > cap(slice) { // if necessary, reallocate
		// allocate double what's needed, for future growth.
		newSlice := make([]string, (n + 1))
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:n]
	slice[n-1] = data
	return slice
}

// Chdir sets the current working directory to the relative path specified.
func Chdir(app string) error {
	var wd, _ = os.Getwd()
	dir, _ := filepath.Abs(filepath.Join(wd, app))
	err := os.Chdir(dir)
	fmt.Println(dir)
	return err
}

// CreateFile creates an empty file at the specified path.
func CreateFile(path string) error {
	fo, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fo.Close()
	return nil
}

// GetFileBytes returns a byte array of the contents of a file.
func GetFileBytes(filename string) []byte {
	file, _ := os.Open(filename)
	defer file.Close()
	stat, _ := file.Stat()
	bs := make([]byte, stat.Size())
	_, _ = file.Read(bs)
	return bs
}

// GetHostIPAddress returns the host IP for accessing etcd and Deis services.
func GetHostIPAddress() string {
	IP := os.Getenv("HOST_IPADDR")
	if IP == "" {
		IP = "172.17.8.100"
	}
	return IP
}

// GetHostOs returns either "darwin" or "ubuntu".
func GetHostOs() string {
	cmd := exec.Command("uname")
	out, _ := cmd.Output()
	if strings.Contains(string(out), "Darwin") {
		return "darwin"
	}
	return "ubuntu"
}

// GetRandomPort returns an unused TCP listen port on the host.
func GetRandomPort() string {
	l, _ := net.Listen("tcp", "127.0.0.1:0") // listen on localhost
	defer l.Close()
	port := l.Addr()
	return strings.Split(port.String(), ":")[1]
}

// NewID returns the first part of a random RFC 4122 UUID
// See http://play.golang.org/p/4FkNSiUDMg
func NewID() string {
	uuid := make([]byte, 16)
	io.ReadFull(rand.Reader, uuid)
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x", uuid[0:4])
}

// Rmdir removes a directory and its contents.
func Rmdir(app string) error {
	var wd, _ = os.Getwd()
	dir, _ := filepath.Abs(filepath.Join(wd, app))
	err := os.RemoveAll(dir)
	fmt.Println(dir)
	return err
}

// RunCommandWithStdoutStderr execs a command and returns its output.
func RunCommandWithStdoutStderr(cmd *exec.Cmd) (bytes.Buffer, bytes.Buffer, error) {
	var stdout, stderr bytes.Buffer
	stderrPipe, err := cmd.StderrPipe()
	stdoutPipe, err := cmd.StdoutPipe()

	cmd.Env = os.Environ()
	if err != nil {
		fmt.Println("error at io pipes")
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("error at command start")
	}

	go func() {
		io.Copy(&stdout, stdoutPipe)
		fmt.Println(stdout.String())
	}()
	go func() {
		io.Copy(&stderr, stderrPipe)
		fmt.Println(stderr.String())
	}()
	time.Sleep(2000 * time.Millisecond)
	err = cmd.Wait()
	if err != nil {
		fmt.Println("error at command wait")
	}
	return stdout, stderr, err
}

func errorOut(err error, t *testing.T, message string) {
	if err != nil {
		t.Fatal(message)
	}
}

func errorOutOnNonNilError(err error, t *testing.T, message string) {
	if err == nil {
		t.Fatalf(message)
	}
}

func getExitCode(err error) (int, error) {
	exitCode := 0
	if exiterr, ok := err.(*exec.ExitError); ok {
		if procExit := exiterr.Sys().(syscall.WaitStatus); ok {
			return procExit.ExitStatus(), nil
		}
	}
	return exitCode, fmt.Errorf("failed to get exit code")
}

func logDone(message string) {
	fmt.Printf("[PASSED]: %s\n", message)
}

func nLines(s string) int {
	return strings.Count(s, "\n")
}

func runCommand(cmd *exec.Cmd) (exitCode int, err error) {
	_, exitCode, err = runCommandWithOutput(cmd)
	return
}

func runCommandWithOutput(cmd *exec.Cmd) (output string, exitCode int, err error) {
	exitCode = 0
	out, err := cmd.CombinedOutput()
	if err != nil {
		var exiterr error
		if exitCode, exiterr = getExitCode(err); exiterr != nil {
			// TODO: Fix this so we check the error's text.
			// we've failed to retrieve exit code, so we set it to 127
			exitCode = 127
		}
	}
	output = string(out)
	return
}

func startCommand(cmd *exec.Cmd) (exitCode int, err error) {
	exitCode = 0
	err = cmd.Start()
	if err != nil {
		var exiterr error
		if exitCode, exiterr = getExitCode(err); exiterr != nil {
			// TODO: Fix this so we check the error's text.
			// we've failed to retrieve exit code, so we set it to 127
			exitCode = 127
		}
	}
	return
}

func stripTrailingCharacters(target string) string {
	target = strings.Trim(target, "\n")
	target = strings.Trim(target, " ")
	return target
}
