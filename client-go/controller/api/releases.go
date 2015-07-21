package api

// Release is the definition of the release object.
type Release struct {
	App     string `json:"app"`
	Build   string `json:"build,omitempty"`
	Config  string `json:"config"`
	Created string `json:"created"`
	Owner   string `json:"owner"`
	Summary string `json:"summary"`
	Updated string `json:"updated"`
	UUID    string `json:"uuid"`
	Version int    `json:"version"`
}

// Releases is the definition of GET /v1/apps/<app id>/releases/.
type Releases struct {
	Count    int       `json:"count"`
	Next     int       `json:"next"`
	Previous int       `json:"previous"`
	Releases []Release `json:"results"`
}

// ReleaseRollback is the defenition of POST /v1/apps/<app id>/releases/.
type ReleaseRollback struct {
	Version int `json:"version"`
}
