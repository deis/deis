package syslog

import "fmt"

type SyslogMessage interface {
	fmt.Stringer
}

// Message defines a syslog message.
type Message struct {
	Msg string
}

func (m *Message) String() string {
	return m.Msg
}
