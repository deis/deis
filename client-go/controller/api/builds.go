package api

// Build is the structure of the build object.
type Build struct {
	App        string            `json:"app"`
	Created    string            `json:"created"`
	Dockerfile string            `json:"dockerfile,omitempty"`
	Image      string            `json:"image,omitempty"`
	Owner      string            `json:"owner"`
	Procfile   map[string]string `json:"procfile"`
	Sha        string            `json:"sha,omitempty"`
	Updated    string            `json:"updated"`
	UUID       string            `json:"uuid"`
}

// Builds is the structure of GET /v1/apps/<app id>/builds/.
type Builds struct {
	Count    int     `json:"count"`
	Next     int     `json:"next"`
	Previous int     `json:"previous"`
	Builds   []Build `json:"results"`
}

// CreateBuildRequest is the structure of POST /v1/apps/<app id>/builds/.
type CreateBuildRequest struct {
	Image    string            `json:"image"`
	Procfile map[string]string `json:"procfile,omitempty"`
}
