package api

// Key is the definition of the key object.
type Key struct {
	Created string `json:"created"`
	ID      string `json:"id"`
	Owner   string `json:"owner"`
	Public  string `json:"public"`
	Updated string `json:"updated"`
	UUID    string `json:"uuid"`
}

type Keys []Key

func (k Keys) Len() int           { return len(k) }
func (k Keys) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }
func (k Keys) Less(i, j int) bool { return k[i].ID < k[j].ID }

// KeyCreateRequest is the definition of POST /v1/keys/.
type KeyCreateRequest struct {
	ID     string `json:"id"`
	Public string `json:"public"`
	Name   string `json:"name,omitempty"`
}
