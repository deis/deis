package net

import (
	"net"
	"strings"
	"time"
)

// WaitForPort wait for successful network connection
func WaitForPort(proto string, ip string, port string, timeout time.Duration) error {
	for {
		con, err := net.DialTimeout(proto, ip+":"+port, timeout)
		if err == nil {
			con.Close()
			break
		}
	}

	return nil
}

func RandomPort() string {
	l, _ := net.Listen("tcp", "127.0.0.1:0") // listen on localhost
	defer l.Close()
	port := l.Addr()
	return strings.Split(port.String(), ":")[1]
}
