package api

// Domain is the structure of the domain object.
type Domain struct {
	App     string `json:"app"`
	Created string `json:"created"`
	Domain  string `json:"domain"`
	Owner   string `json:"owner"`
	Updated string `json:"updated"`
}

// Domains is the structure of GET /v1/app/<app id>/domains/.
type Domains struct {
	Count    int      `json:"count"`
	Next     int      `json:"next"`
	Previous int      `json:"previous"`
	Domains  []Domain `json:"results"`
}

// DomainCreateRequest is the structure of POST /v1/app/<app id>/domains/.
type DomainCreateRequest struct {
	Domain string `json:"domain"`
}
