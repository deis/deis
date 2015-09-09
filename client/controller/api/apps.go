package api

// App is the definition of the app object.
type App struct {
	Created string `json:"created"`
	ID      string `json:"id"`
	Owner   string `json:"owner"`
	Updated string `json:"updated"`
	URL     string `json:"url"`
	UUID    string `json:"uuid"`
}

// AppCreateRequest is the definition of POST /v1/apps/.
type AppCreateRequest struct {
	ID string `json:"id,omitempty"`
}

// AppUpdateRequest is the definition of POST /v1/apps/<app id>/.
type AppUpdateRequest struct {
	Owner string `json:"owner,omitempty"`
}

// AppRunRequest is the definition of POST /v1/apps/<app id>/run.
type AppRunRequest struct {
	Command string `json:"command"`
}

// AppRunResponse is the definition of /v1/apps/<app id>/run.
type AppRunResponse struct {
	Output     string `json:"output"`
	ReturnCode int    `json:"rc"`
}
