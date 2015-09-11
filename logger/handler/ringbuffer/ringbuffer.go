package ringbuffer

import (
	"container/ring"
	"fmt"
	"sync"
)

type ringBuffer struct {
	ring  *ring.Ring
	mutex *sync.Mutex
}

func newRingBuffer(size int) *ringBuffer {
	return &ringBuffer{
		ring:  ring.New(size),
		mutex: &sync.Mutex{},
	}
}

func (rb *ringBuffer) write(message string) {
	rb.mutex.Lock()
	defer rb.mutex.Unlock()
	rb.ring = rb.ring.Next()
	rb.ring.Value = message
}

func (rb *ringBuffer) read(lines int) []string {
	if lines <= 0 {
		return []string{}
	}
	rb.mutex.Lock()
	defer rb.mutex.Unlock()
	var start *ring.Ring
	if lines < rb.ring.Len() {
		start = rb.ring.Move(-1 * (lines - 1))
	} else {
		start = rb.ring.Next()
	}
	data := make([]string, 0, lines)
	start.Do(func(line interface{}) {
		if line == nil || lines <= 0 {
			return
		}
		lines -= 1
		data = append(data, line.(string))
	})
	return data
}

type handler struct {
	bufferSize  int
	ringBuffers map[string]*ringBuffer
	mutex       sync.Mutex
}

func NewHandler(bufferSize int) (*handler, error) {
	if bufferSize <= 0 {
		return nil, fmt.Errorf("Invalid ringBuffer size: %d", bufferSize)
	}
	return &handler{bufferSize: bufferSize, ringBuffers: make(map[string]*ringBuffer)}, nil
}

// Write adds a log message to to an app-specific ringBuffer
func (h *handler) Write(app string, message string) error {
	rb, ok := h.ringBuffers[app]
	if !ok {
		h.mutex.Lock()
		defer h.mutex.Unlock()
		rb, ok = h.ringBuffers[app]
		if !ok {
			rb = newRingBuffer(h.bufferSize)
			h.ringBuffers[app] = rb
		}
	}
	rb.write(message)
	return nil
}

// Read retrieves a specified number of log lines from an app-specific ringBuffer
func (h *handler) Read(app string, lines int) ([]string, error) {
	rb, ok := h.ringBuffers[app]
	if ok {
		return rb.read(lines), nil
	}
	return nil, fmt.Errorf("Could not find logs for '%s'", app)
}
