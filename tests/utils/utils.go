// Package utils contains commonly useful functions from Deis testing.
package utils

import (
	"bufio"
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
)

// BuildTag returns the $BUILD_TAG environment variable or `git rev-parse` output.
func BuildTag() string {
	var tag string
	tag = os.Getenv("BUILD_TAG")
	if tag == "" {
		out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
		if err != nil {
			tag = "NOTAG"
		}
		tag = "git-" + string(out)
	}
	return strings.TrimSpace(tag)
}

// ImagePrefix returns the $IMAGE_PREFIX environment variable or `deis/`
func ImagePrefix() string {
	var prefix string
	prefix = os.Getenv("IMAGE_PREFIX")
	if prefix != "" {
		return prefix
	}
	return "deis/"
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

// HostAddress returns the host IP for accessing etcd and Deis services.
func HostAddress() string {
	IP := os.Getenv("HOST_IPADDR")
	if IP == "" {
		IP = "172.17.8.100"
	}
	return IP
}

// Hostname returns the hostname of the machine running the container, *not* the local machine
// We infer the hostname because we don't necessarily know how to log in.
func Hostname() string {
	switch HostAddress() {
	case "172.17.8.100":
		return "deis-1"
	case "172.17.8.101":
		return "deis-2"
	case "172.17.8.102":
		return "deis-3"
	case "172.21.12.100":
		return "docker-registry"
	default:
		return "boot2docker"
	}
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

// RandomPort returns an unused TCP listen port on the host.
func RandomPort() string {
	l, _ := net.Listen("tcp", "127.0.0.1:0") // listen on localhost
	defer l.Close()
	port := l.Addr()
	return strings.Split(port.String(), ":")[1]
}

// Rmdir removes a directory and its contents.
func Rmdir(app string) error {
	var wd, _ = os.Getwd()
	dir, _ := filepath.Abs(filepath.Join(wd, app))
	err := os.RemoveAll(dir)
	fmt.Println(dir)
	return err
}

// streamOutput from a source to a destination buffer while also printing
func streamOutput(src io.Reader, dst *bytes.Buffer, out io.Writer) error {

	s := bufio.NewReader(src)

	for {
		var line []byte
		line, err := s.ReadSlice('\n')
		if err == io.EOF && len(line) == 0 {
			break // done
		}
		if err == io.EOF {
			return fmt.Errorf("Improper termination: %v", line)
		}
		if err != nil {
			return err
		}

		// append to the buffer
		dst.Write(line)

		// write to stdout/stderr also
		out.Write(line)
	}

	return nil
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
		streamOutput(stdoutPipe, &stdout, os.Stdout)
	}()
	go func() {
		streamOutput(stderrPipe, &stderr, os.Stderr)
	}()
	err = cmd.Wait()
	if err != nil {
		fmt.Println("error at command wait")
	}
	return stdout, stderr, err
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
