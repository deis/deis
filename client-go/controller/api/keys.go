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

// Keys is the definition of GET /v1/keys/.
type Keys struct {
	Count    int   `json:"count"`
	Next     int   `json:"next"`
	Previous int   `json:"previous"`
	Keys     []Key `json:"results"`
}

// KeyCreateRequest is the definition of POST /v1/keys/.
type KeyCreateRequest struct {
	ID     string `json:"id"`
	Public string `json:"public"`
	Name   string `json:"name,omitempty"`
}
