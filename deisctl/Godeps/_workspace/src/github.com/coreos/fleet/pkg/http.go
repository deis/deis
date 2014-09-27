package pkg

import (
	"net/http"

	"github.com/coreos/fleet/log"
)

type LoggingHTTPTransport struct {
	http.Transport
}

func (lt *LoggingHTTPTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	log.V(1).Infof("HTTP %s %s", req.Method, req.URL.String())
	resp, err = lt.Transport.RoundTrip(req)
	if err == nil {
		log.V(1).Infof("HTTP %s %s %s", req.Method, req.URL.String(), resp.Status)
	}
	return
}
