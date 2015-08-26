package logspout

import (
	"net"
)

const (
	// This can be tweaked to allow longer log lines
	writeBufferSize = 1024 * 1024
)

type deisTransport int

func newDeisTransport() *deisTransport {
	return new(deisTransport)
}

func (_ *deisTransport) Dial(addr string, options map[string]string) (net.Conn, error) {
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}
	err = conn.SetWriteBuffer(writeBufferSize)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
