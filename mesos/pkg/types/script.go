package types

// Script struct to specify a script.
type Script struct {
	Name    string
	Params  map[string]string
	Content func(string) ([]byte, error)
}
