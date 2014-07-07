package syslog

import (
    "testing"
    "time"
)

type TestLogger struct {}

func (tl TestLogger) Print(...interface{}) {}
func (tl TestLogger) Printf(format string, v ...interface{}) {}
func (tl TestLogger) Println(...interface{}) {}

func TestNewFileHandler(t *testing.T) {
    fh := NewFileHandler("", 1, func (m *Message) bool {return true}, true)
    if fh == nil {
        t.Errorf("expected filehandler, got nil")
    }
}

func TestSetLogger(t *testing.T) {
    tl := TestLogger{}
    fh := NewFileHandler("", 1, func (m *Message) bool {return true}, true)
    fh.SetLogger(tl)
    if fh.l != tl {
        t.Errorf("expected the logger to be set")
    }
}

func TestHandle(t *testing.T) {
    fh := NewFileHandler("/tmp/test", 1, func (m *Message) bool {return true}, true)
    handle := fh.Handle(&Message{time.Now(), nil, 0, 0, time.Now(), "localhost", "test", "message", "", ""})
    if handle == nil {
        t.Errorf("expected a handle, got nil")
    }
}
