package syslogish

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"

	"github.com/deis/deis/logger/drain"
	"github.com/deis/deis/logger/storage"
)

const queueSize = 500

var appRegex *regexp.Regexp

func init() {
	appRegex = regexp.MustCompile(`^.* ([-_a-z0-9]+)\[[a-z0-9-_\.]+\].*`)
}

// Server implements a UDP-based "syslog-like" server.  Like syslog, as described by RFC 3164, it
// expects that each packet contains a single log message and that, conversely, log messages are
// encapsulated in their entirety by a single packet, however, no attempt is made to parse the
// messages received or validate that they conform to the specification.
type Server struct {
	conn           net.PacketConn
	listening      bool
	storageQueue   chan string
	storageAdapter storage.Adapter
	drainageQueue  chan string
	drain          drain.LogDrain
	adapterMutex   sync.RWMutex
	drainMutex     sync.RWMutex
}

// NewServer returns a pointer to a new Server instance.
func NewServer(bindHost string, bindPort int) (*Server, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", bindHost, bindPort))
	if err != nil {
		return nil, err
	}
	c, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		conn:          c,
		storageQueue:  make(chan string, queueSize),
		drainageQueue: make(chan string, queueSize),
	}, nil
}

// SetStorageAdapter permits a server's underlying storage.Adapter to be reconfigured (replaced)
// at runtime.
func (s *Server) SetStorageAdapter(storageAdapter storage.Adapter) {
	// Get an exclusive lock before updating the internal pointer to the storage adapter.  Other
	// goroutines holding read locks might depend on that pointer as it currently exists.
	s.adapterMutex.Lock()
	defer s.adapterMutex.Unlock()
	s.storageAdapter = storageAdapter
}

// SetDrain permits a server's underlying drain.LogDrain to be reconfigured (replaced) at runtime.
func (s *Server) SetDrain(drain drain.LogDrain) {
	// Get an exclusive lock before updating the internal pointer to the log drain.  Other
	// goroutines holding read locks might depend on that pointer as it currently exists.
	s.drainMutex.Lock()
	defer s.drainMutex.Unlock()
	s.drain = drain
}

// Listen starts the server's main loop.
func (s *Server) Listen() {
	// Should only ever be called once
	if !s.listening {
		s.listening = true
		go s.receive()
		go s.processStorage()
		go s.processDrainage()
		log.Println("syslogish server running")
	}
}

func (s *Server) receive() {
	// Make buffer the same size as the max for a UDP packet
	buf := make([]byte, 65535)
	for {
		n, _, err := s.conn.ReadFrom(buf)
		if err != nil {
			log.Fatal("syslogish server read error", err)
		}
		message := strings.TrimSuffix(string(buf[:n]), "\n")
		select {
		case s.storageQueue <- message:
		default:
		}
	}
}

func (s *Server) processStorage() {
	for message := range s.storageQueue {
		app, err := getAppName(message)
		if err != nil {
			log.Println(err)
			return
		}
		// Get a read lock to ensure the storage adapater pointer can't be nilled by the configurer
		// in the time between we check if it's nil and the time we invoke .Write() upon it.
		s.adapterMutex.RLock()
		// DONT'T defer unlocking... defered statements are executed when the function returns, but
		// we are inside an infinite loop here.  If we defer, we would never release the lock.
		// Instead, release it manually below.
		if s.storageAdapter != nil {
			s.storageAdapter.Write(app, message)
			// We don't bother trapping errors here, so failed writes to storage are silent.  This is by
			// design.  If we sent a log message to STDOUT in response to the failure, deis-logspout
			// would read it and forward it back to deis-logger, which would fail again to write to
			// storage and spawn ANOTHER log message.  The effect would be an infinite loop of
			// unstoreable log messages that would nevertheless fill up journal logs and eventually
			// overake the disk.
			//
			// Treating this as a fatal event would cause the deis-logger unit to restart-- sending
			// even more log messages to STDOUT.  The overall effect would be the same as described
			// above with the added disadvantages of flapping.
		}
		s.adapterMutex.RUnlock()
		// Add the message to the drainage queue.  This allows the storage loop to continue right
		// away instead of waiting while the message is sent to an external service-- since that
		// could be a bottleneck and error prone depending on rate limiting, network congestion, etc.
		select {
		case s.drainageQueue <- message:
		default:
		}
	}
}

func (s *Server) processDrainage() {
	for message := range s.drainageQueue {
		// Get a read lock to ensure the drain pointer can't be nilled by the configurer in the time
		// between we check if it's nil and the time we invoke .Send() upon it.
		s.drainMutex.RLock()
		// DONT'T defer unlocking... defered statements are executed when the function returns, but
		// we are inside an infinite loop here.  If we defer, we would never release the lock.
		// Instead, release it manually below.
		if s.drain != nil {
			s.drain.Send(message)
			// We don't bother trapping errors here, so failed sends to the drain are silent.  This is
			// by design.  If we sent a log message to STDOUT in response to the failure, deis-logspout
			// would read it and forward it back to deis-logger, which would fail again to send to the
			// drain and spawn ANOTHER log message.  The effect would be an infinite loop of undrainable
			// log messages that would nevertheless fill up journal logs and eventually overake the disk.
			//
			// Treating this as a fatal event would cause the deis-logger unit to restart-- sending
			// even more log messages to STDOUT.  The overall effect would be the same as described
			// above with the added disadvantages of flapping.
		}
		s.drainMutex.RUnlock()
	}
}

func getAppName(message string) (string, error) {
	match := appRegex.FindStringSubmatch(message)
	if match == nil {
		return "", fmt.Errorf("Could not find app name in message: %s", message)
	}
	return match[1], nil
}

// ReadLogs returns a specified number of log lines (if available) for a specified app by
// delegating to the server's underlying storage.Adapter.
func (s *Server) ReadLogs(app string, lines int) ([]string, error) {
	// Get a read lock to ensure the storage adapater pointer can't be updated by another
	// goroutine in the time between we check if it's nil and the time we invoke .Read() upon
	// it.
	s.adapterMutex.RLock()
	defer s.adapterMutex.RUnlock()
	if s.storageAdapter == nil {
		return nil, fmt.Errorf("Could not find logs for '%s'.  No storage adapter specified.", app)
	}
	return s.storageAdapter.Read(app, lines)
}

// DestroyLogs deletes all logs for a specified app by delegating to the server's underlying
// storage.Adapter.
func (s *Server) DestroyLogs(app string) error {
	// Get a read lock to ensure the storage adapater pointer can't be updated by another
	// goroutine in the time between we check if it's nil and the time we invoke .Destroy() upon
	// it.
	s.adapterMutex.RLock()
	defer s.adapterMutex.RUnlock()
	if s.storageAdapter == nil {
		return fmt.Errorf("Could not destroy logs for '%s'.  No storage adapter specified.", app)
	}
	return s.storageAdapter.Destroy(app)
}

// ReopenLogs delegate to the server's underlying storage.Adapter to, if applicable, refresh
// references to underlying storage mechanisms.  This is useful, for instance, to ensure logging
// continues smoothly after log rotation when file-based storage is in use.
func (s *Server) ReopenLogs() error {
	// Get a read lock to ensure the storage adapater pointer can't be updated by another
	// goroutine in the time between we check if it's nil and the time we invoke .Reopen() upon
	// it.
	s.adapterMutex.RLock()
	defer s.adapterMutex.RUnlock()
	if s.storageAdapter == nil {
		return errors.New("Could not reopen logs.  No storage adapter specified.")
	}
	return s.storageAdapter.Reopen()
}
