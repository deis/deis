package main

import (
	"os"
	"os/exec"
	"fmt"
	"net"
	"io"
)

const (
	CONN_HOST = "0.0.0.0"
	CONN_TYPE = "tcp"
)

// Handles incoming connection
func handleRequest(conn net.Conn, command string) {
	cmd := exec.Command("bash", "-c", command)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	cmd.Start()

	connInCh := make(chan []byte)
	connInErroCh := make(chan error)
	stdoutCh := make(chan []byte)
	stdoutErroCh := make(chan error)

	// Setup connection channel
	go readInput(conn, connInCh, connInErroCh)
	// Setup stdout channel
	go readInput(stdout, stdoutCh, stdoutErroCh)

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
		}
	}

	// Determine exit code
	err = cmd.Wait();
	statusCode := 0;
	if err != nil {
		statusCode = 1;
	}
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
	port := os.Getenv("PORT")

	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + port)
	if err != nil {
			fmt.Println("Error listening:", err.Error())
			os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Listening on " + CONN_HOST + ":" + port)
	// Listen for an incoming connection.
	conn, err := l.Accept()
	if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
	}

	// Handle connection
	handleRequest(conn, command)
}
