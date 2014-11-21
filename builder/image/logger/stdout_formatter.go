package logger

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
)

type StdOutFormatter struct {
}

func (f *StdOutFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "[%s] - %s\n", strings.ToUpper(entry.Level.String()), entry.Message)
	return b.Bytes(), nil
}
