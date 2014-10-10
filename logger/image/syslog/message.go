package syslog

import (
	"fmt"
	"strings"
)

type SyslogMessage interface {
	fmt.Stringer
}

// Message defines a syslog message.
type Message struct {
	Msg string
}

func (m *Message) String() string {
	return strings.TrimSuffix(m.Msg, "\n")
}
