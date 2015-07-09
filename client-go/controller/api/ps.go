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

// Processes defines the structure of processes.
type Processes struct {
	Count     int       `json:"count"`
	Next      int       `json:"next"`
	Previous  int       `json:"previous"`
	Processes []Process `json:"results"`
}
