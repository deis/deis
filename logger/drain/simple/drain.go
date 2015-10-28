package simple

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"sync"
	"time"
)

// For efficiency, we reuse connections for a while (instead of dialing every time).  However,
// there are two compelling reasons to redial periodically:
//
//  1. We don't want DNS changes on the remote end of the drain to go unnoticed for too long.
//
//  2. If the drain is using TCP, the underlying TCP stack can potentially take a very long time
//     waiting for acks and retrying send for packets that haven't been acked.  This creates a
//     large window where packets can be spewed into the ether (without any warning) before the
//     problem is detected.  By redialing periodically, we create the opportunity for a failed TCP
//     handshake-- which tells us sooner that something is wrong.
//
// For efficiency we want the refresh interval to be high.  For resiliency, we want it to be low.
// One minute has been arbitrarily selected as a sensible balance of these two concerns.
const connRefreshInterval = 1 * time.Minute

// This determines how many failed dial attempts are required before the drain is muted.
const maxFailedConns = 5

// This determines how much time we're willing to spend dialing.
const dialTimeout = 10 * time.Second

// This is how long the drain is muted for after repeated connection failures.
const mutePeriod = 5 * time.Minute

type logDrain struct {
	proto string
	uri   string
	conn  net.Conn
	muted bool
	mutex sync.Mutex
}

// NewDrain returns a pointer to a new instance of a drain.LogDrain
func NewDrain(drainURL string) (*logDrain, error) {
	u, err := url.Parse(drainURL)
	if err != nil {
		return nil, err
	}
	var proto string
	if u.Scheme == "udp" || u.Scheme == "syslog" {
		proto = "udp"
	} else if u.Scheme == "tcp" {
		proto = "tcp"
	} else {
		return nil, fmt.Errorf("Invalid drain url scheme: %s", u.Scheme)
	}
	return &logDrain{proto: proto, uri: u.Host + u.Path}, nil
}

// Send forwards the provided log message to an external destination
func (d *logDrain) Send(message string) error {
	if d.muted {
		return nil
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	conn, err := d.getConnection(false)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(conn, message)
	if err != nil {
		// Try again with a new connection in case the issue was a broken pipe
		conn, err = d.getConnection(true)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(conn, message)
		if err != nil {
			return err
		}
	}
	return nil
}

// getConnection returns a usable connection, often without needing to redial, but still
// redialing when advised.
func (d *logDrain) getConnection(forceNew bool) (net.Conn, error) {
	// If we have a connection, it's not old, and we're not focing a new one...
	if d.conn != nil && !forceNew {
		// then return the existing connection
		return d.conn, nil
	}
	// If ANY of those conditions weren't met, it's time for a new connection.
	// If we have an existing one, close it and nil it out, too for good measure.
	if d.conn != nil {
		if err := d.conn.Close(); err != nil {
			log.Println("drain: Error closing connection.  Drain may be leaking connections.", err)
		}
		d.conn = nil
	}
	// Try a few times...
	var err error
	for attempt := 1; attempt <= maxFailedConns; attempt++ {
		d.conn, err = net.DialTimeout(d.proto, d.uri, dialTimeout)
		if err == nil {
			// We got our connection...
			// Make it good for only so long.  See comment above on connRefreshInterval.
			err = d.conn.SetWriteDeadline(time.Now().Add(connRefreshInterval))
			if err != nil {
				return nil, err
			}
			// Break out of the loop and return
			return d.conn, nil
		}
	}
	// Multiple attempts to dial have failed.  Whatever the problem is, we shouldn't expect that
	// it will resolve itself quickly.
	log.Printf("drain: Experienced %d consecutive failed connection attempts; muting drain for %s.", maxFailedConns, mutePeriod)
	// Immediately "mute" the drain.  This will prevent us from wasting resources repeatedly dialing
	// and failing while the message queue gets backed up.  This will give the network a break and
	// allow us to empty the queue.
	d.muted = true
	// Unmute the drain when the mute interval has elapsed
	go func() {
		time.Sleep(mutePeriod)
		d.muted = false
	}()
	// Return the error from the last failed connection attempt
	return nil, err
}
