package api

// Cert is the definition of the cert object.
// Some fields are omtempty because they are only
// returned when creating or getting a cert.
type Cert struct {
	Updated string `json:"updated,omitempty"`
	Created string `json:"created,omitempty"`
	Name    string `json:"common_name"`
	Expires string `json:"expires"`
	Owner   string `json:"owner,omitempty"`
	ID      int    `json:"id,omitempty"`
}

// CertCreateRequest is the definition of POST /v1/certs/.
type CertCreateRequest struct {
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
	Name        string `json:"common_name,omitempty"`
}
