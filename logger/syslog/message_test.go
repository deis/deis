package syslog

import (
    "net"
    "testing"
    "time"
)

func TestMessageNetSrc(t *testing.T) {
    tcpAddress, err := net.ResolveTCPAddr("tcp", "localhost:1234")

    if err != nil {
        t.Errorf("could not resolve TCP address")
    }

    m := &Message{time.Now(), tcpAddress, 0, 0, time.Now(), "", "", "", "", ""}
    if m.NetSrc() != "127.0.0.1" {
        t.Errorf("expected 127.0.0.1, got ", m.NetSrc())
    }

    udpAddress, err := net.ResolveUDPAddr("udp", "localhost:1234")

    if err != nil {
        t.Errorf("could not resolve UDP address")
    }

    m.Source = udpAddress
    if m.NetSrc() != "127.0.0.1" {
        t.Errorf("expected 127.0.0.1, got ", m.NetSrc())
    }

    unixAddress, err := net.ResolveUnixAddr("unix", "/tmp/str")

    if err != nil {
        t.Errorf("could not resolve unix address")
    }

    m.Source = unixAddress
    if m.NetSrc() != "/tmp/str" {
        t.Errorf("expected /tmp/str, got ", m.NetSrc())
    }

    unknownAddress, err := net.ResolveIPAddr("ip4", "localhost")

    if err != nil {
        t.Errorf("could not resolve unknown address")
    }

    m.Source = unknownAddress
    if m.NetSrc() != "127.0.0.1" {
        t.Errorf("expected 127.0.0.1, got ", m.NetSrc())
    }
}

func TestMessageFormat(t *testing.T) {
    tcpAddress, err := net.ResolveTCPAddr("tcp", "localhost:1234")

    if err != nil {
        t.Errorf("could not resolve TCP address")
    }

    m := &Message{
        time.Now(),
        tcpAddress,
        0,
        0,
        time.Now(),
        "localhost",
        "TEST",
        "hello world",
        "",
        "",
    }

    timeLayout := "2006-01-02 15:04:05"
    expectedOutput := m.Time.Format(timeLayout) + " localhost hello world"
    if m.String() != expectedOutput {
        t.Errorf("expected '" + expectedOutput + "', got '" + m.String() + "'.")
    }
}
