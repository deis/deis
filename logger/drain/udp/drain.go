package udp

import (
	"fmt"
	"net"
	"net/url"
)

type logDrain struct {
	uri string
}

// NewDrain returns a pointer to a new instance of a UDP-based drain.LogDrain.
func NewDrain(drainURL string) (*logDrain, error) {
	u, err := url.Parse(drainURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "udp" && u.Scheme != "syslog" {
		return nil, fmt.Errorf("Invalid drain url scheme: %s", u.Scheme)
	}
	return &logDrain{uri: u.Host + u.Path}, nil
}

// Send forwards the provided log message to an external destination using UDP for transport.
func (d *logDrain) Send(message string) error {
	conn, err := net.Dial("udp", d.uri)
	if err != nil {
		return fmt.Errorf("Error dialing log drain at %s over udp", d.uri)
	}
	defer conn.Close()
	fmt.Fprintln(conn, message)
	return nil
}
