package logger

import (
	"os"

	"github.com/Sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	Log.Formatter = new(StdOutFormatter)
	// use the env LOG to set the log level
	logLevel := os.Getenv("LOG")
	if logLevel != "" {
		if level, err := logrus.ParseLevel(logLevel); err == nil {
			Log.Level = level
		}
	}
}
