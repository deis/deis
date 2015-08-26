package api

// ConfigSet is the definition of POST /v1/apps/<app id>/config/.
type ConfigSet struct {
	Values map[string]string `json:"values"`
}

// ConfigUnset is the definition of POST /v1/apps/<app id>/config/.
type ConfigUnset struct {
	Values map[string]interface{} `json:"values"`
}

// Config is the structure of an app's config.
type Config struct {
	Owner   string                 `json:"owner,omitempty"`
	App     string                 `json:"app,omitempty"`
	Values  map[string]interface{} `json:"values,omitempty"`
	Memory  map[string]interface{} `json:"memory,omitempty"`
	CPU     map[string]interface{} `json:"cpu,omitempty"`
	Tags    map[string]interface{} `json:"tags,omitempty"`
	Created string                 `json:"created,omitempty"`
	Updated string                 `json:"updated,omitempty"`
	UUID    string                 `json:"uuid,omitempty"`
}
