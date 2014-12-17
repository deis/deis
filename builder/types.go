package builder

import (
	"github.com/deis/deis/pkg/time"
)

// ProcessType represents the key/value mappings of a process type to a process inside
// a Heroku Procfile.
//
// See https://devcenter.heroku.com/articles/procfile
type ProcessType map[string]string

// ConfigHook represents a repository from which to extract the configuration and user to use.
type ConfigHook struct {
	ReceiveUser string `json:"receive_user"`
	ReceiveRepo string `json:"receive_repo"`
}

// BuildHook represents a controller's build-hook object.
type BuildHook struct {
	Sha         string      `json:"sha"`
	ReceiveUser string      `json:"receive_user"`
	ReceiveRepo string      `json:"receive_repo"`
	Image       string      `json:"image"`
	Procfile    ProcessType `json:"procfile"`
	Dockerfile  string      `json:"dockerfile"`
}

// BuildHookResponse represents a controller's build-hook response object.
type BuildHookResponse struct {
	Release map[string]int `json:"release"`
	Domains []string       `json:"domains"`
}

// Config represents a Deis application's configuration.
type Config struct {
	Owner   string                 `json:"owner"`
	App     string                 `json:"app"`
	Values  map[string]interface{} `json:"values"`
	Memory  map[string]string      `json:"memory"`
	CPU     map[string]int         `json:"cpu"`
	Tags    map[string]string      `json:"tags"`
	UUID    string                 `json:"uuid"`
	Created time.Time              `json:"created"`
	Updated time.Time              `json:"updated"`
}
