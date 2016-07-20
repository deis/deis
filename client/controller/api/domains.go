package api

// Domain is the structure of the domain object.
type Domain struct {
	App     string `json:"app"`
	Created string `json:"created"`
	Domain  string `json:"domain"`
	Owner   string `json:"owner"`
	Updated string `json:"updated"`
}

type Domains []Domain

func (d Domains) Len() int           { return len(d) }
func (d Domains) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d Domains) Less(i, j int) bool { return d[i].Domain < d[j].Domain }

// DomainCreateRequest is the structure of POST /v1/app/<app id>/domains/.
type DomainCreateRequest struct {
	Domain string `json:"domain"`
}
