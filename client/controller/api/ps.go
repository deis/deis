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

type Processes []Process

func (p Processes) Len() int           { return len(p) }
func (p Processes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Processes) Less(i, j int) bool { return p[i].Num < p[j].Num }

type ProcessType struct {
	Type      string
	Processes Processes
}

type ProcessTypes []ProcessType

func (p ProcessTypes) Len() int           { return len(p) }
func (p ProcessTypes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ProcessTypes) Less(i, j int) bool { return p[i].Type < p[j].Type }
