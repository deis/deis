package drain

import (
	"fmt"
	"strings"

	"github.com/deis/deis/logger/drain/simple"
)

// NewDrain returns a pointer to an appropriate implementation of the LogDrain interface, as
// determined by the drainURL it is passed.
func NewDrain(drainURL string) (LogDrain, error) {
	if drainURL == "" {
		// nil means no drain-- which is valid
		return nil, nil
	}
	// Any of these three can use the same drain implementation
	if strings.HasPrefix(drainURL, "udp://") || strings.HasPrefix(drainURL, "syslog://") || strings.HasPrefix(drainURL, "tcp://") {
		drain, err := simple.NewDrain(drainURL)
		if err != nil {
			return nil, err
		}
		return drain, nil
	}
	// TODO: Add more drain implementations-- TLS over TCP and HTTP/S
	return nil, fmt.Errorf("Cannot construct a drain for URL: '%s'", drainURL)
}
