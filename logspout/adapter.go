package logspout

import (
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gliderlabs/logspout/router"
)

const deisDatetimeFormat string = "2006-01-02T15:04:05MST"

type deisLogAdapter struct {
	conn  net.Conn
	route *router.Route
}

func newDeisLogAdapter(route *router.Route) (router.LogAdapter, error) {
	transport, found := router.AdapterTransports.Lookup(route.AdapterTransport("deis-udp"))
	if !found {
		return nil, errors.New("bad transport: " + route.Adapter)
	}
	conn, err := transport.Dial(route.Address, route.Options)
	if err != nil {
		return nil, err
	}
	return &deisLogAdapter{
		route: route,
		conn:  conn,
	}, nil
}

func (a *deisLogAdapter) Stream(logstream chan *router.Message) {
	for message := range logstream {
		m := &DeisLogMessage{message}
		_, err := a.conn.Write(m.Render())
		if err != nil {
			log.Println("deis:", err)
			if reflect.TypeOf(a.conn).String() != "*net.UDPConn" {
				return
			}
		}
	}
}

type DeisLogMessage struct {
	*router.Message
}

func (m *DeisLogMessage) Render() []byte {
	tag, pid := getLogName(m.Container.Name)
	return []byte(fmt.Sprintf("%s %s[%s]: %s", time.Now().Format(deisDatetimeFormat), tag, pid, m.Data))
}

// getLogName returns a custom tag and PID for containers that
// match Deis' specific application name format. Otherwise,
// it returns the original name and 1 as the PID.
func getLogName(name string) (string, string) {
	// example regex that should match: go_v2.web.1
	match := getMatch(`^\/?([a-z0-9-]+)_(v[0-9]+)\.([a-z-_]+\.[0-9]+)$`, name)
	if match != nil {
		return match[1], match[3]
	}
	match = getMatch(`^\/?k8s_([a-z0-9-]+)-[a-z]+\.[\da-f]+_[a-z0-9-]+-([a-z]+-[\da-z]*)_`, name)
	if match != nil {
		return match[1], match[2]
	}
	return strings.TrimPrefix(name, "/"), "1"
}

func getMatch(regex string, name string) []string {
	r := regexp.MustCompile(regex)
	match := r.FindStringSubmatch(name)
	return match
}
