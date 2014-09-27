package pkg

import (
	"time"
)

func ExpBackoff(last time.Duration, max time.Duration) (next time.Duration) {
	next = last * 2
	if next > max {
		next = max
	}
	return
}
