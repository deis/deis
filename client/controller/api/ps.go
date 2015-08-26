package api

// Process defines the structure of a process.
type Process struct {
	Owner   string `json:"owner"`
	App     string `json:"app"`
	Release string `json:"release"`
	Created string `json:"created"`
	Updated string `json:"updated"`
	UUID    string `json:"uuid"`
	Type    string `json:"type"`
	Num     int    `json:"num"`
	State   string `json:"state"`
}
