package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ziutek/syslog"
)

type handler struct {
	*syslog.BaseHandler
}

func newHandler() *handler {
	h := handler{syslog.NewBaseHandler(5, func(m *syslog.Message) bool { return true }, false)}
	go h.mainLoop()
	return &h
}

func (h *handler) mainLoop() {
	for {
		m := h.Get()
		if m == nil {
			break
		}
		fmt.Println(m)
	}
	h.End()
}

func main() {
	flag.Parse()
	s := syslog.NewServer()
	s.AddHandler(newHandler())
	s.Listen(flag.Arg(0))

	// Wait for terminating signal
	sc := make(chan os.Signal, 2)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT)
	<-sc

	fmt.Println("Shutdown the server...")
	s.Shutdown()
}
