// Syslog server library. It is based on RFC 3164 so it doesn't parse properly
// packets with new header format (described in RFC 5424).
package syslog

import (
	"bytes"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Server struct {
	conns    []net.PacketConn
	handlers []Handler
	shutdown bool
	l        FatalLogger
}

//  NewServer creates idle server
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

func (s *Server) passToHandlers(m *Message) {
	for _, h := range s.handlers {
		m = h.Handle(m)
		if m == nil {
			break
		}
	}
}

func (s *Server) receiver(c net.PacketConn) {
	//q := (chan<- Message)(s.q)
	buf := make([]byte, 1024)
	for {
		n, addr, err := c.ReadFrom(buf)
		if err != nil {
			if !s.shutdown {
				s.l.Fatalln("Read error:", err)
			}
			return
		}
		pkt := buf[:n]

		m := new(Message)
		m.Source = addr
		m.Time = time.Now()

		// Parse priority (if exists)
		prio := 13 // default priority
		hasPrio := false
		if pkt[0] == '<' {
			n = 1 + bytes.IndexByte(pkt[1:], '>')
			if n > 1 && n < 5 {
				p, err := strconv.Atoi(string(pkt[1:n]))
				if err == nil && p >= 0 {
					hasPrio = true
					prio = p
					pkt = pkt[n+1:]
				}
			}
		}
		m.Severity = Severity(prio & 0x07)
		m.Facility = Facility(prio >> 3)

		// Parse header (if exists)
		if hasPrio && len(pkt) >= 16 && pkt[15] == ' ' {
			// Get timestamp
			layout := "Jan _2 15:04:05"
			ts, err := time.Parse(layout, string(pkt[:15]))
			if err == nil && !ts.IsZero() {
				// Get hostname
				n = 16 + bytes.IndexByte(pkt[16:], ' ')
				if n != 15 {
					m.Timestamp = ts
					m.Hostname = string(pkt[16:n])
					pkt = pkt[n+1:]
				}
			}
			// TODO: check for version an new format of header as
			// described in RFC 5424.
		}

		// Parse msg part
		msg := string(bytes.TrimRightFunc(pkt, isNulCrLf))
		n = strings.IndexFunc(msg, isNotAlnum)
		if n != -1 {
			m.Tag = msg[:n]
			m.Content = msg[n:]
		} else {
			m.Content = msg
		}
		msg = strings.TrimFunc(msg, unicode.IsSpace)
		n = strings.IndexFunc(msg, unicode.IsSpace)
		if n != -1 {
			m.Tag1 = msg[:n]
			m.Content1 = strings.TrimLeftFunc(msg[n+1:], unicode.IsSpace)
		} else {
			m.Content1 = msg
		}

		s.passToHandlers(m)
	}
}
