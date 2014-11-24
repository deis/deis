package commons

import (
	"os"

	"github.com/deis/deis/controller/logger"
)

// Getopt return the value of and environment variable or a default
func Getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		logger.Log.Debugf("returning default value \"%s\" for key \"%s\"", dfault, name)
		value = dfault
	}
	return value
}
