package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"net"
	"fmt"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.net/websocket"

	"github.com/go-martini/martini"
)

// global session
var session *Session

type Session struct {
	sync.Mutex
	Conns map[string]*Conn
}

type Conn struct {
	Conn net.Conn
}

func (s *Session) Get(name string) (*Conn, error) {
	s.Lock()
	defer s.Unlock()
	conn, found := s.Conns[name]
	if !found {
		return nil, errors.New("could not find session")
	}
	return conn, nil
}

func (s *Session) Create(name string, conn *Conn) error {
	if _, ok := s.Conns[name]; ok == true {
		return errors.New("session already exists")
	}
	s.Lock()
	defer s.Unlock()
	log.Println("creating session", name)
	s.Conns[name] = conn
	return nil
}

func (s *Session) Delete(name string) {
	s.Lock()
	defer s.Unlock()
	log.Println("deleting session", name)
	delete(s.Conns, name)
}

func NewSession() *Session {
	return &Session{sync.Mutex{}, make(map[string]*Conn, 0)}
}

func CreateSession(req *http.Request, res http.ResponseWriter) (int, string) {
	// read in the form values
	if err := req.ParseForm(); err != nil {
		return 400, "could not parse form"
	}

	host := req.Form.Get("host")
	if host == "" {
		return 400, "no host supplied"
	}

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create a new session ID
	sessionID := uuid.New()

	session.Create(sessionID, &Conn{
		Conn: conn,
	})
	return 200, sessionID
}

func SendMessage(ws *websocket.Conn) {
	req := ws.Request()
	sessionID := strings.Split(req.URL.Path, "/")[2]
	conn, err := session.Get(sessionID)
	if err != nil {
		ws.Write([]byte("could not retrieve session"))
	}

	connInCh := make(chan []byte)
	connInErroCh := make(chan error)
	stdoutCh := make(chan []byte)
	stdoutErroCh := make(chan error)

	// Setup connection channel
	go readInput(conn.Conn, connInCh, connInErroCh)
	// Setup stdout channel
	go readInput(ws, stdoutCh, stdoutErroCh)

	// Pipe in- and outputs
	loop:
	for {
		select {
			// Received data from connection
			case data := <-connInCh:
				ws.Write(data)
			// Received an error from connection
			case err := <-connInErroCh:
				fmt.Println("Error conn:", err.Error())
				break loop;
			// Received data from stdout
			case data := <-stdoutCh:
				conn.Conn.Write(data)
			// Received an error from stdout
			case err := <-stdoutErroCh:
				fmt.Println("Error stdout:", err.Error())
				break loop;
		}
	}

	ws.Close()
	conn.Conn.Close()
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

func DeleteSession(params martini.Params) (int, string) {
	session.Delete(params["id"])
	return 204, "NO CONTENT"
}

func server() *martini.Martini {
	session = NewSession()
	m := martini.New()
	r := martini.NewRouter()
	r.Post(`/sessions`, CreateSession)
	r.Get(`/sessions/:id`, websocket.Handler(SendMessage).ServeHTTP)
	r.Delete(`/sessions/:id`, DeleteSession)
	// Add the router action
	m.Action(r.Handle)
	return m
}

func NewServer() *martini.Martini {
	m := server()
	// Setup middleware
	m.Use(martini.Static("."))
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	return m
}

func NewTestServer() *martini.Martini {
	return server()
}

func main() {
	m := NewServer()
	m.RunOnAddr("0.0.0.0:"+os.Getenv("PORT"))
}
