package auth

import (
	"crypto/sha256"
	"crypto/tls"
	"net/http"

	"github.com/coreos/updateservicectl/Godeps/_workspace/src/github.com/tent/hawk-go"
)

var DefaultHawkHasher = sha256.New

type HawkRoundTripper struct {
	User          string
	Token         string
	SkipSSLVerify bool
}

func (t *HawkRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	creds := &hawk.Credentials{
		ID:   t.User,
		Key:  t.Token,
		Hash: DefaultHawkHasher,
	}

	auth := hawk.NewRequestAuth(req, creds, 0)
	req.Header.Set("Authorization", auth.RequestHeader())

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: t.SkipSSLVerify},
	}
	return transport.RoundTrip(req)
}
