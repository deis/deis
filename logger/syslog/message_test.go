package syslog

import (
	"testing"
	"time"
)

func TestMessageFormat(t *testing.T) {

	m := &Message{time.Now(), 34, time.Now(), "localhost", "test", "hello world"}

	timeLayout := "2006-01-02 15:04:05"
	expectedOutput := "<34>" + m.Time.Format(timeLayout) + " test: hello world"
	if m.String() != expectedOutput {
		t.Errorf("expected '" + expectedOutput + "', got '" + m.String() + "'.")
	}
}
