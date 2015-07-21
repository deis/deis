package log

import (
	"os"

	"github.com/Sirupsen/logrus"
)

// Logger embed logrus Logger struct
type Logger struct {
	logrus.Logger
}

// New create a new logger using the StdOutFormatter and the level
// specified in the env variable LOG_LEVEL
func New() *Logger {
	log := &Logger{}

	log.Out = os.Stdout
	log.Formatter = new(StdOutFormatter)

	logLevel := os.Getenv("LOG_LEVEL")
	log.SetLevel(logLevel)

	return log
}

// SetLevel change the level of the logger
func (log *Logger) SetLevel(logLevel string) {
	if logLevel != "" {
		if level, err := logrus.ParseLevel(logLevel); err == nil {
			log.Level = level
		}
	}
}
