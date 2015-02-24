package logger

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
)

// StdOutFormatter formats log messages from the router component.
type StdOutFormatter struct {
}

// Format rewrites a log entry for stdout as a byte array.
func (f *StdOutFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "[%s] - %s\n", strings.ToUpper(entry.Level.String()), entry.Message)
	return b.Bytes(), nil
}
