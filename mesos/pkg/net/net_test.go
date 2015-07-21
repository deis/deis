package net

import (
	"net"
	"strconv"
	"testing"
	"time"
)

func TestListenTCP(t *testing.T) {
	port, err := RandomPort("tcp")
	if err != nil {
		t.Fatal(err)
	}

	listeningPort, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	defer listeningPort.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = WaitForPort("tcp", "127.0.0.1", port, time.Second)
	if err != nil {
		t.Fatal(err)
	}
}

// TODO: fix
// func TestListenUDP(t *testing.T) {
// 	port, err := RandomPort("udp")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1"+strconv.Itoa(port))
// 	listeningPort, err := net.ListenUDP("udp", addr)
// 	defer listeningPort.Close()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	err = WaitForPort("udp", "127.0.0.1", port, time.Second)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
