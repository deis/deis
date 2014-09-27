package pkg

import (
	"testing"
	"time"
)

func TestExpBackoff(t *testing.T) {
	tests := []struct {
		last time.Duration
		max  time.Duration
		next time.Duration
	}{
		{1 * time.Second, 10 * time.Second, 2 * time.Second},
		{8 * time.Second, 10 * time.Second, 10 * time.Second},
		{10 * time.Second, 10 * time.Second, 10 * time.Second},
		{20 * time.Second, 10 * time.Second, 10 * time.Second},
	}

	for i, tt := range tests {
		next := ExpBackoff(tt.last, tt.max)
		if next != tt.next {
			t.Errorf("case %d: last=%v, max=%v, next=%v; got next=%v", i, tt.last, tt.max, tt.next, next)
		}
	}
}
