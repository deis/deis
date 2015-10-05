package drain

import (
	"fmt"
	"strings"

	"github.com/deis/deis/logger/drain/tcp"
	"github.com/deis/deis/logger/drain/udp"
)

// NewDrain returns a pointer to an appropriate implementation of the LogDrain interface, as
// determined by the drainURL it is passed.
func NewDrain(drainURL string) (LogDrain, error) {
	if drainURL == "" {
		// nil means no drain-- which is valid
		return nil, nil
	}
	if strings.HasPrefix(drainURL, "udp://") || strings.HasPrefix(drainURL, "syslog://") {
		drain, err := udp.NewDrain(drainURL)
		if err != nil {
			return nil, err
		}
		return drain, nil
	}
	if strings.HasPrefix(drainURL, "tcp://") {
		drain, err := tcp.NewDrain(drainURL)
		if err != nil {
			return nil, err
		}
		return drain, nil
	}
	return nil, fmt.Errorf("Cannot construct a drain for URL: '%s'", drainURL)
}
