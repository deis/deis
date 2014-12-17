package syslog

// Handler handles syslog messages
type Handler interface {
	// Handle should return Message (mayby modified) for future processing by
	// other handlers or return nil. If Handle is called with nil message it
	// should complete all remaining work and properly shutdown before return.
	Handle(*Message) *Message
}

// BaseHandler is desigend for simplify the creation of real handlers. It
// implements Handler interface using nonblocking queuing of messages and
// simple message filtering.
type BaseHandler struct {
	queue  chan *Message
	end    chan struct{}
	filter func(*Message) bool
	ft     bool
}

// NewBaseHandler creates BaseHandler using specified filter. If filter is nil
// or if it returns true messages are passed to BaseHandler internal queue
// (of qlen length). If filter returns false or ft is true messages are returned
// to server for future processing by other handlers.
func NewBaseHandler(qlen int, filter func(*Message) bool, ft bool) *BaseHandler {
	return &BaseHandler{
		queue:  make(chan *Message, qlen),
		end:    make(chan struct{}),
		filter: filter,
		ft:     ft,
	}
}

// Handle inserts m in an internal queue. It immediately returns even if
// queue is full. If m == nil it closes queue and waits for End method call
// before return.
func (h *BaseHandler) Handle(m *Message) *Message {
	if m == nil {
		close(h.queue) // signal that ther is no more messages for processing
		<-h.end        // wait for handler shutdown
		return nil
	}
	if h.filter != nil && !h.filter(m) {
		// m doesn't match the filter
		return m
	}
	// Try queue m
	select {
	case h.queue <- m:
	default:
	}
	if h.ft {
		return m
	}
	return nil
}

// Get returns first message from internal queue. It waits for message if queue
// is empty. It returns nil if there is no more messages to process and handler
// should shutdown.
func (h *BaseHandler) Get() *Message {
	m, ok := <-h.queue
	if ok {
		return m
	}
	return nil
}

// Queue returns BaseHandler internal queue as read-only channel. You can use
// it directly, especially if your handler need to select from multiple channels
// or have to work without blocking. You need to check if this channel is closed by
// sender and properly shutdown in this case.
func (h *BaseHandler) Queue() <-chan *Message {
	return h.queue
}

// End signals the server that handler properly shutdown. You need to call End
// only if Get has returned nil before.
func (h *BaseHandler) End() {
	close(h.end)
}
