// Package syslog implements a syslog server library. It is based on RFC 3164,
// as such it does not properly parse packets with an RFC 5424 header format.
package syslog

import (
	"log"
	"net"
	"os"
	"strings"
	"unicode"
)

// Server is the wrapper for a syslog server.
type Server struct {
	conns    []net.PacketConn
	handlers []Handler
	shutdown bool
	l        FatalLogger
}

// NewServer creates an idle server.
func NewServer() *Server {
	return &Server{l: log.New(os.Stderr, "", log.LstdFlags)}
}

// SetLogger sets logger for server errors. A running server is rather quiet and
// logs only fatal errors using FatalLogger interface. By default standard Go
// logger is used so errors are writen to stderr and after that whole
// application is halted. Using SetLogger you can change this behavior (log
// erross elsewhere and don't halt whole application).
func (s *Server) SetLogger(l FatalLogger) {
	s.l = l
}

// AddHandler adds h to internal ordered list of handlers
func (s *Server) AddHandler(h Handler) {
	s.handlers = append(s.handlers, h)
}

// Listen starts gorutine that receives syslog messages on specified address.
// addr can be a path (for unix domain sockets) or host:port (for UDP).
func (s *Server) Listen(addr string) error {
	var c net.PacketConn
	if strings.IndexRune(addr, ':') != -1 {
		a, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			return err
		}
		c, err = net.ListenUDP("udp", a)
		if err != nil {
			return err
		}
	} else {
		a, err := net.ResolveUnixAddr("unixgram", addr)
		if err != nil {
			return err
		}
		c, err = net.ListenUnixgram("unixgram", a)
		if err != nil {
			return err
		}
	}
	s.conns = append(s.conns, c)
	go s.receiver(c)
	return nil
}

// Shutdown stops server.
func (s *Server) Shutdown() {
	s.shutdown = true
	for _, c := range s.conns {
		err := c.Close()
		if err != nil {
			s.l.Fatalln(err)
		}
	}
	s.passToHandlers(nil)
	s.conns = nil
	s.handlers = nil
}

func isNotAlnum(r rune) bool {
	return !(unicode.IsLetter(r) || unicode.IsNumber(r))
}

func isNulCrLf(r rune) bool {
	return r == 0 || r == '\r' || r == '\n'
}

func (s *Server) passToHandlers(m SyslogMessage) {
	for _, h := range s.handlers {
		m = h.Handle(m)
		if m == nil {
			break
		}
	}
}

func (s *Server) receiver(c net.PacketConn) {
	// make packet buffer the same size as logspout
	buf := make([]byte, 1048576)
	for {
		n, _, err := c.ReadFrom(buf)
		if err != nil {
			if !s.shutdown {
				s.l.Fatalln("Read error:", err)
			}
			return
		}
		// pass along the incoming syslog message
		s.passToHandlers(&Message{string(buf[:n])})
	}
}
