package main

import (
	"os"
	"os/exec"
	"log"
	"net"
	"io"
	"time"
	"fmt"
)

const (
	connHost = "0.0.0.0"
	connType = "tcp"
	connPort = "3333"
	timeout time.Duration = 10 * time.Second
	etcdTTL time.Duration = timeout * 2
	connectionTimeout time.Duration = 10 * time.Minute
)

// Handles incoming connection
func handleRequest(conn net.Conn, command string) {
	log.Print("Executing: ", command)

	cmd := exec.Command("bash", "-c", command)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	cmd.Start()

	connInCh := make(chan []byte)
	connInErroCh := make(chan error)
	stdoutCh := make(chan []byte)
	stdoutErroCh := make(chan error)
	stderrCh := make(chan []byte)
	stderrErroCh := make(chan error)

	// Setup connection channel
	go readInput(conn, connInCh, connInErroCh)
	// Setup stdout channel
	go readInput(stdout, stdoutCh, stdoutErroCh)
	// Setup stderr channel
	go readInput(stderr, stderrCh, stderrErroCh)

	// Pipe in- and outputs
	loop:
	for {
		select {
			// Received data from connection
			case data := <-connInCh:
				stdin.Write(data)
			// Received an error from connection
			case <-connInErroCh:
				break loop;
			// Received data from stdout
			case data := <-stdoutCh:
				conn.Write(data)
			// Received an error from stdout
			case <-stdoutErroCh:
				break loop;
			// Received data from stdout
			case data := <-stderrCh:
				conn.Write(data)
			// Received an error from stdout
			case <-stderrErroCh:
				break loop;
		}
	}

	// Determine exit code
	statusCode := 0;
	done := make(chan error, 1)
	go func() {
	    done <- cmd.Wait()
	}()
	select {
	    case <-time.After(10 * time.Second):
	        if err := cmd.Process.Kill(); err != nil {
	            log.Print("Failed to kill: ", err)
	        }
	        <-done
					log.Print("Process killed.")
	        statusCode = 1;
	    case err := <-done:
          if err != nil {
						statusCode = 1;
					}
	}

	log.Print("Exited with: ", statusCode)
	statusResponse := []byte(fmt.Sprintf("%v", statusCode));
	conn.Write(statusResponse)

	conn.Close()
}

func readInput(r io.Reader, ch chan []byte, eCh chan error) {
	for {
		data := make([]byte, 512)
		_, err := r.Read(data)
		if err != nil {
			eCh<- err
			return
		}
		ch<- data
	}
}

func main() {
	command := os.Getenv("COMMAND")

	// Listen for incoming connections.
	tcpAddr, err := net.ResolveTCPAddr("tcp", connHost + ":" + connPort)
	if err != nil {
		log.Fatal("Error resolving:", err.Error())
	}
	l, err := net.ListenTCP(connType, tcpAddr)
	if err != nil {
		log.Fatal("Error listening:", err.Error())
	}
	// Timeout after 10 minutes without connection
	l.SetDeadline(time.Now().Add(connectionTimeout))
	// Close the listener when the application closes.
	defer l.Close()

	log.Print("Listening on " + connHost + ":" + connPort)

	// Listen for an incoming connection.
	conn, err := l.Accept()
	if err != nil {
		log.Fatal("Error accepting: ", err.Error())
	}

	// Handle connection
	handleRequest(conn, command)
}
