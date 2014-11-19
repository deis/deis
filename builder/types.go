package builder

import (
	"time"
)

const (
	DEIS_DATETIME string = "2006-01-02T15:04:05MST"
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
	Dockerfile  bool        `json:"dockerfile"`
}

// BuildHookResponse represents a controller's build-hook response object.
type BuildHookResponse struct {
	Release map[string]int `json:"release"`
	Domains []string       `json:"domains"`
}

// Config represents a Deis application's configuration.
type Config struct {
	Owner   string            `json:"owner"`
	App     string            `json:"app"`
	Values  map[string]string `json:"values"`
	Memory  map[string]string `json:"memory"`
	CPU     map[string]string `json:"cpu"`
	Tags    map[string]string `json:"tags"`
	UUID    string            `json:"uuid"`
	Created DeisTime          `json:"created"`
	Updated DeisTime          `json:"updated"`
}

// DeisTime represents the standard datetime format used across the platform.
type DeisTime struct {
	time.Time
}

func (t *DeisTime) format() string {
	return t.Time.Format(DEIS_DATETIME)
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in Deis' datetime format.
func (t *DeisTime) MarshalJSON() ([]byte, error) {
	return []byte(t.Time.Format(`"` + DEIS_DATETIME + `"`)), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in Deis' datetime format.
func (t *DeisTime) UnmarshalText(data []byte) (err error) {
	tt, err := time.Parse(DEIS_DATETIME, string(data))
	*t = DeisTime{tt}
	return
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in Deis' datetime format.
func (t *DeisTime) UnmarshalJSON(data []byte) (err error) {
	// Fractional seconds are handled implicitly by Parse.
	tt, err := time.Parse(`"`+DEIS_DATETIME+`"`, string(data))
	*t = DeisTime{tt}
	return
}
