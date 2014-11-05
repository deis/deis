// Package utils contains commonly useful functions from Deisctl

package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/deis/deis/deisctl/constant"
	uuid "github.com/satori/go.uuid"
)

// NewUUID returns a new V4-style unique identifier.
func NewUUID() string {
	u1 := uuid.NewV4()
	s1 := fmt.Sprintf("%s", u1)
	return strings.Split(s1, "-")[0]
}

func getetcdClient() *etcd.Client {
	machines := []string{"http://127.0.0.1:4001/"}
	return etcd.NewClient(machines)
}

// GetKey returns the value of an etcd key, or of an environment variable if etcd didn't have
// a value.
func GetKey(dir, key, perm string) string {
	c := getetcdClient()
	result, err := c.Get(dir+key, false, false)
	if err != nil || result.Node.Value == "" {
		return os.Getenv(perm)
	}
	return result.Node.Value
}

// GetClientID returns the CoreOS Machine ID, or an unknown UUID string.
func GetClientID() string {
	machineID := GetMachineID("/")
	if machineID == "" {
		return fmt.Sprintf("{unknown-" + NewUUID() + "}")
	}
	return machineID
}

// GetMachineID returns the CoreOS Machine ID.
func GetMachineID(root string) string {
	fullPath := filepath.Join(root, constant.MachineID)
	id, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(id))
}

// GetVersion returns the package version from a text file resource.
func GetVersion() string {
	id, err := ioutil.ReadFile(constant.Version)
	if err != nil {
		return "0.0.0"
	}
	return strings.TrimSpace(string(id))
}

// Extract expands a .tar archive file into the specified directory.
func Extract(file, dir string) (err error) {
	var wd, _ = os.Getwd()
	_ = os.Chdir(dir)
	cmdl := exec.Command("tar", "-C", "/", "-xvf", file)
	if _, _, err := RunCommandWithStdoutStderr(cmdl); err != nil {
		fmt.Printf("Failed:\n%v", err)
		return err
	}
	_ = os.Chdir(wd)
	return nil
}

// PutVersion updates the package version to a local text file resource.
func PutVersion(version string) error {
	return ioutil.WriteFile(constant.Version, []byte(version), 0644)
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

// Execute runs the given script in a shell.
func Execute(script string) error {
	cmdl := exec.Command("sh", "-c", script)
	if _, _, err := RunCommandWithStdoutStderr(cmdl); err != nil {
		fmt.Println("(Error )")
		return err
	}
	return nil
}

// DeisIfy returns a pretty-printed deis logo along with the corresponding message
func DeisIfy(message string) string {
	circle := "\033[31m●"
	square := "\033[32m■"
	triangle := "\033[34m▴"
	reset := "\033[0m"
	title := reset + message

	return fmt.Sprintf("%s %s %s\n%s %s %s %s\n%s %s %s%s\n",
		circle, triangle, square,
		square, circle, triangle, title,
		triangle, square, circle, reset)
}
