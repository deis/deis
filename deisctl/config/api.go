package config

// Client interface used for configuration
type Client interface {
	Get(string) (string, error)
	Set(string) (string, error)
}
