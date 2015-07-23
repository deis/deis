package api

// PermsAppResponse is the definition of GET /v1/apps/<app id>/perms/.
type PermsAppResponse struct {
	Users []string `json:"users"`
}

// PermsAdminResponse is the definition of GET /v1/admin/perms/.
type PermsAdminResponse struct {
	Count    int `json:"count"`
	Next     int `json:"next"`
	Previous int `json:"previous"`
	Users    []struct {
		Username string `json:"username"`
	} `json:"results"`
}

// PermsRequest is the definition of a requst on /perms/.
type PermsRequest struct {
	Username string `json:"username"`
}
