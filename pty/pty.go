package main

import (
	"github.com/gorilla/websocket"
	"github.com/kr/pty"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
	"fmt"
)

const (
	port = "3333"
	closeTimeout = 5*time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1,
	WriteBufferSize: 1,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsPty struct {
	Cmd *exec.Cmd // pty builds on os.exec
	Pty *os.File  // a pty is simply an os.File
}

func (wp *wsPty) Start() {
	var err error
	command := os.Getenv("COMMAND")
	wp.Cmd = exec.Command("bash", "-c", command)
	wp.Pty, err = pty.Start(wp.Cmd)
	if err != nil {
		log.Fatalf("Failed to start command: %s\n", err)
	}
}

func (wp *wsPty) Stop() {
	wp.Pty.Close()
}

func ptyHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Websocket upgrade failed: %s\n", err)
	}
	defer conn.Close()

	wp := wsPty{}
	wp.Start()

	connInCh := make(chan []byte)
	connInErroCh := make(chan error)
	stdoutCh := make(chan []byte)
	stdoutErroCh := make(chan error)

	// copy everything from the pty master to the websocket
	go func() {
		for {
			buf := make([]byte, 512)
			_, err := wp.Pty.Read(buf)
			if err != nil {
				stdoutErroCh<- err
				return
			}
			stdoutCh<- buf
		}
	}()

	// read from the web socket, copying to the pty master
	go func() {
		for {
			mt, payload, err := conn.ReadMessage()
			if err != nil {
				connInErroCh<- err
				return
			}

			switch mt {
			case websocket.TextMessage:
				connInCh<- payload
			default:
				log.Printf("Invalid message type %d\n", mt)
				return
			}
		}
	}()

	// Pipe in- and outputs
	loop:
	for {
		select {
			// Received data from connection
			case data := <-connInCh:
				wp.Pty.Write(data)
			// Received an error from connection
			case <-connInErroCh:
				break loop;
			// Received data from stdout
			case data := <-stdoutCh:
				conn.WriteMessage(websocket.TextMessage, data)
			// Received an error from stdout
			case <-stdoutErroCh:
				break loop;
		}
	}

	// Determine exit code
	statusCode := 0;
	done := make(chan error, 1)
	go func() {
	    done <- wp.Cmd.Wait()
	}()
	select {
	    case <-time.After(10 * time.Second):
	        if err := wp.Cmd.Process.Kill(); err != nil {
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

	closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure,
																							 fmt.Sprintf("%v", statusCode))

	conn.WriteControl(websocket.CloseMessage,
										closeMessage,
										time.Now().Add(closeTimeout))
	wp.Stop()
	os.Exit(0)
}


func main() {
	http.HandleFunc("/", ptyHandler)
	log.Print("Listening on 0.0.0.0:" + port)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
