package drain

// LogDrain is an interface for pluggable components that ship logs to a remote destination.
type LogDrain interface {
	Send(string) error
}
