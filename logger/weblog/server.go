package weblog

import (
	"fmt"
	"log"
	"net/http"

	"github.com/deis/deis/logger/syslogish"
)

// Server implements a simple HTTP server that handles GET and DELETE requests for application
// logs.  These actions are accomplished by delegating to a syslogish.Server, which will broker
// communication between its underlying storage.Adapter and this weblog server.
type Server struct {
	listening bool
	bindHost  string
	bindPort  int
	handler   *requestHandler
}

// NewServer returns a pointer to a new Server instance.
func NewServer(bindHost string, bindPort int, syslogishServer *syslogish.Server) (*Server, error) {
	return &Server{
		bindHost: bindHost,
		bindPort: bindPort,
		handler:  &requestHandler{syslogishServer: syslogishServer},
	}, nil
}

// Listen starts the server's main loop.
func (s *Server) Listen() {
	// Should only ever be called once
	if !s.listening {
		s.listening = true
		go s.listen()
		log.Println("weblog server running")
	}
}

func (s *Server) listen() {
	http.Handle("/", s.handler)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.bindHost, s.bindPort), nil); err != nil {
		log.Fatal("weblog server stopped", err)
	}
}
