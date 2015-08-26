package api

// PermsAppResponse is the definition of GET /v1/apps/<app id>/perms/.
type PermsAppResponse struct {
	Users []string `json:"users"`
}

// PermsRequest is the definition of a requst on /perms/.
type PermsRequest struct {
	Username string `json:"username"`
}
