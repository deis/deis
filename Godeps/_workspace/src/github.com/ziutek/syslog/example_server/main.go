package main

import (
	"fmt"
	"github.com/ziutek/syslog"
	"os"
	"os/signal"
	"syscall"
)

type handler struct {
	// To simplify implementation of our handler we embed helper
	// syslog.BaseHandler struct.
	*syslog.BaseHandler
}

// Simple fiter for named/bind messages which can be used with BaseHandler
func filter(m *syslog.Message) bool {
	return m.Tag == "named" || m.Tag == "bind"
}

func newHandler() *handler {
	h := handler{syslog.NewBaseHandler(5, filter, false)}
	go h.mainLoop() // BaseHandler needs some gorutine that reads from its queue
	return &h
}

// mainLoop reads from BaseHandler queue using h.Get and logs messages to stdout
func (h *handler) mainLoop() {
	for {
		m := h.Get()
		if m == nil {
			break
		}
		fmt.Println(m)
	}
	fmt.Println("Exit handler")
	h.End()
}

func main() {
	// Create a server with one handler and run one listen gorutine
	s := syslog.NewServer()
	s.AddHandler(newHandler())
	s.Listen("0.0.0.0:1514")

	// Wait for terminating signal
	sc := make(chan os.Signal, 2)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT)
	<-sc

	// Shutdown the server
	fmt.Println("Shutdown the server...")
	s.Shutdown()
	fmt.Println("Server is down")
}
