package net

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"time"
)

// InterfaceIPAddress is used to know the interface and ip address in the sytem
type InterfaceIPAddress struct {
	Iface string
	IP    string
}

// WaitForPort wait for successful network connection
func WaitForPort(proto string, ip string, port int, timeout time.Duration) error {
	for {
		con, err := net.DialTimeout(proto, ip+":"+strconv.Itoa(port), timeout)
		if err == nil {
			con.Close()
			break
		}
	}

	return nil
}

// RandomPort return a random not used TCP port
func RandomPort(proto string) (int, error) {
	switch proto {
	case "tcp":
		l, _ := net.Listen(proto, "127.0.0.1:0")
		defer l.Close()
		port := l.Addr()
		lPort, _ := strconv.Atoi(strings.Split(port.String(), ":")[1])
		return lPort, nil
	case "udp":
		addr, _ := net.ResolveUDPAddr(proto, "127.0.0.1:0")
		l, _ := net.ListenUDP(proto, addr)
		defer l.Close()
		return addr.Port, nil
	default:
		return -1, errors.New("invalid protocol")
	}
}

// GetNetworkInterfaces return the list of
// network interfaces and IP address
func GetNetworkInterfaces() []InterfaceIPAddress {
	result := []InterfaceIPAddress{}

	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if addrs, err := inter.Addrs(); err == nil {
			for _, addr := range addrs {
				result = append(result, InterfaceIPAddress{inter.Name, addr.String()})
			}
		}
	}

	return result
}

// ParseIP parses s as an IP address
func ParseIP(s string) net.IP {
	return net.ParseIP(s)
}
