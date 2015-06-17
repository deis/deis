package config

// Client interface used for configuration
type Client interface {
	Get(string) (string, error)
	Set(string, string) (string, error)
	Delete(string) error
}
