// Package utils contains commonly useful functions from Deis testing.

package utils

import (
	_ "bufio"
	"bytes"
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

	"github.com/satori/go.uuid"
)

// NewUuid returns a new V4-style unique identifier.
func NewUuid() string {
	u1 := uuid.NewV4()
	s1 := fmt.Sprintf("%s", u1)
	return strings.Split(s1, "-")[0]
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

func ListFiles(dir string) ([]string, error) {
	files, err := filepath.Glob(dir)
	return files, err
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

// Chdir sets the current working directory to the relative path specified.
func Chdir(app string) error {
	var wd, _ = os.Getwd()
	dir, _ := filepath.Abs(filepath.Join(wd, app))
	err := os.Chdir(dir)
	fmt.Println(dir)
	return err
}

func Extract(file, dir string) {
	Chdir(dir)
	cmdl := exec.Command("tar", "-xvf", file)
	if _, _, err := RunCommandWithStdoutStderr(cmdl); err != nil {
		fmt.Printf("Failed:\n%v", err)
	} else {
		fmt.Println("ok")
	}
	Chdir("/home/core/")
}

// Rmdir removes a directory and its contents.
func Rmdir(app string) error {
	var wd, _ = os.Getwd()
	dir, _ := filepath.Abs(filepath.Join(wd, app))
	err := os.RemoveAll(dir)
	fmt.Println(dir)
	return err
}

// GetUserDetails returns sections of a UUID.
func GetUserDetails() (string, string) {
	u1 := uuid.NewV4()
	s1 := fmt.Sprintf("%s", u1)
	return strings.Split(s1, "-")[0], strings.Split(s1, "-")[1]
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

// GetHostIPAddress returns the host IP for accessing etcd and Deis services.
func GetHostIPAddress() string {
	IP := os.Getenv("HOST_IPADDR")
	if IP == "" {
		IP = "172.17.8.100"
	}
	return IP
}

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

// GetRandomPort returns an unused TCP listen port on the host.
func GetRandomPort() string {
	l, _ := net.Listen("tcp", "127.0.0.1:0") // listen on localhost
	defer l.Close()
	port := l.Addr()
	return strings.Split(port.String(), ":")[1]
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

func runCommand(cmd *exec.Cmd) (exitCode int, err error) {
	exitCode = 0
	err = cmd.Run()
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

func logDone(message string) {
	fmt.Printf("[PASSED]: %s\n", message)
}

func stripTrailingCharacters(target string) string {
	target = strings.Trim(target, "\n")
	target = strings.Trim(target, " ")
	return target
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

func nLines(s string) int {
	return strings.Count(s, "\n")
}

//func deis(bash string , arg string ,  cmd string )
