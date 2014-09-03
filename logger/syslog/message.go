package syslog

import (
	"fmt"
	"log/syslog"
	"time"
)

// Message defines an RFC 3164 syslog message.
type Message struct {
	Time      time.Time
	Priority  syslog.Priority
	Timestamp time.Time
	Hostname  string
	Tag       string
	Content   string
}

// String returns the Message in a string format. This satisfies the fmt.Stringer
// interface.
func (m *Message) String() string {
	timeLayout := "2006-01-02 15:04:05"
	return fmt.Sprintf(
		"<%d>%s %s: %s",
		m.Priority,
		m.Time.Format(timeLayout),
		m.Tag,
		m.Content,
	)
}
