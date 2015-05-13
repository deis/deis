package log

import (
	"os"

	"github.com/Sirupsen/logrus"
)

// New create a new logger using the StdOutFormatter and the level
// specified in the env variable LOG_LEVEL
func New() *logrus.Logger {
	log := logrus.New()
	log.Formatter = new(StdOutFormatter)

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		if level, err := logrus.ParseLevel(logLevel); err == nil {
			log.Level = level
		}
	}

	return log
}
