package types

type Script struct {
	Name    string
	Params  map[string]string
	Content func(string) ([]byte, error)
}
