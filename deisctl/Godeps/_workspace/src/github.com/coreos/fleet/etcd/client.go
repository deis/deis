package etcd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/coreos/fleet/third_party/github.com/golang/glog"
)

const (
	defaultEndpoint = "http://localhost:4001"
	redirectMax     = 10
	actionTimeout   = time.Second
)

type Client interface {
	Do(Action) (*Result, error)
	Wait(Action, <-chan bool) (*Result, error)
}

func NewClient(endpoints []string, transport http.Transport) (*client, error) {
	if len(endpoints) == 0 {
		endpoints = []string{defaultEndpoint}
	}

	parsed := make([]url.URL, len(endpoints))
	for i, ep := range endpoints {
		u, err := url.Parse(ep)
		if err != nil {
			return nil, err
		}

		setDefaultPath(u)

		if err = filterURL(u); err != nil {
			return nil, err
		}

		parsed[i] = *u
	}

	return &client{parsed, &transport}, nil
}

// setDefaultPath will set the Path attribute of the provided
// url.URL to / if no Path is set
func setDefaultPath(u *url.URL) {
	// Set a default path
	if u.Path == "" {
		u.Path = "/"
	}
}

// filterURL raises an error if the provided url.URL has any
// questionable attributes
func filterURL(u *url.URL) error {
	if u.Scheme != "http" {
		return fmt.Errorf("unable to use endpoint scheme %s, http only", u.Scheme)
	}

	if u.Path != "/" {
		return fmt.Errorf("unable to use endpoint with non-root path: %s", u)
	}

	if len(u.Query()) > 0 {
		return fmt.Errorf("unable to use endpoint with query parameters: %s", u)
	}

	if len(u.Opaque) > 0 {
		return fmt.Errorf("malformed endpoint: %s", u)
	}

	if u.User != nil {
		return fmt.Errorf("unable to use endpoint with user info: %s", u)
	}

	if len(u.Fragment) > 0 {
		return fmt.Errorf("unable to use endpoint with fragment: %s", u)
	}

	return nil
}

type client struct {
	endpoints []url.URL
	transport *http.Transport
}

// a requestFunc must never return a nil *http.Response and a nil error together
type requestFunc func(*http.Request, <-chan bool) (*http.Response, []byte, error)

// Make a single http request, draining the body on success. If the request
// fails, an error is returned. If the provided channel is ever closed, the
// in-flight request will be cancelled.
func (c *client) requestHTTP(req *http.Request, cancel <-chan bool) (resp *http.Response, body []byte, err error) {
	respchan := make(chan *http.Response, 1)
	errchan := make(chan error, 1)

	go func() {
		resp, err := c.transport.RoundTrip(req)
		if resp == nil && err == nil {
			err = errors.New("nil error and nil response")
		}

		if err != nil {
			errchan <- err
			if resp != nil {
				resp.Body.Close()
			}
			return
		}

		respchan <- resp
	}()

	select {
	case resp = <-respchan:
		defer resp.Body.Close()
	case err = <-errchan:
		return nil, nil, err
	case <-cancel:
		go c.transport.CancelRequest(req)
		return nil, nil, errors.New("cancelled")
	}

	body, err = ioutil.ReadAll(resp.Body)
	return resp, body, err
}

// Attempt to get a usable Result for the provided Action.
// - this call will block until the provided channel is closed
// - requests are attempted against all configured endpoints
// - exponential backoff is used before reattempting resolution
//   of the given Action against the set of endpoints
// - up to 10 redirects are followed per endpoint per attempt
// If the provided channel is closed before a Result can be
// retrieved, a nil object is returned.
func (c *client) resolve(act Action, rf requestFunc, cancel <-chan bool) (*Result, error) {
	requests := func() (res *Result, err error) {
		for eIndex := 0; eIndex < len(c.endpoints); eIndex++ {
			endpoint := c.endpoints[eIndex]
			ar := newActionResolver(act, &endpoint, rf)
			res, err = ar.Resolve(cancel)
			if res != nil || err != nil {
				break
			}

			select {
			case <-cancel:
				return
			default:
			}
		}

		return
	}

	backoff := func(fn func() (*Result, error)) (res *Result, err error) {
		sleep := 100 * time.Millisecond
		for {
			res, err = fn()
			if res != nil || err != nil {
				break
			}

			select {
			case <-cancel:
				return nil, errors.New("cancelled")
			default:
			}

			log.Errorf("Unable to get result for %v, retrying in %v", act, sleep)

			select {
			case <-cancel:
				return nil, errors.New("cancelled")
			case <-time.After(sleep):
			}

			sleep = sleep * 2
			if sleep > time.Second {
				sleep = time.Second
			}
		}
		return
	}

	return backoff(requests)
}

// Make any necessary HTTP requests to resolve the given Action, returning
// a Result if one can be acquired. This function call will wait 10s before
// aborting any in-flight requests and returning an error.
func (c *client) Do(act Action) (*Result, error) {
	cancel := make(chan bool)
	go func() {
		<-time.After(actionTimeout)
		close(cancel)
	}()

	return c.resolve(act, c.requestHTTP, cancel)
}

// Make any necessary HTTP requests to resolve the given Action, returning
// a Result if one can be acquired. If the provided channel is ever closed,
// all in-flight HTTP requests will be aborted and an error will be returned.
func (c *client) Wait(act Action, cancel <-chan bool) (*Result, error) {
	return c.resolve(act, c.requestHTTP, cancel)
}

var (
	handlers = map[int]func(*http.Response, []byte) (*Result, error){
		http.StatusOK:                 unmarshalSuccessfulResponse,
		http.StatusCreated:            unmarshalSuccessfulResponse,
		http.StatusNotFound:           unmarshalFailedResponse,
		http.StatusPreconditionFailed: unmarshalFailedResponse,
		http.StatusBadRequest:         unmarshalFailedResponse,
	}
)

type actionResolver struct {
	action      Action
	endpoint    *url.URL
	requestFunc requestFunc

	redirectCount int
}

func newActionResolver(act Action, ep *url.URL, rf requestFunc) *actionResolver {
	return &actionResolver{action: act, endpoint: ep, requestFunc: rf}
}

// Resolve attempts to yield a result from the configured action and endpoint. If a usable
// Result or error was not attained, nil values are returned.
func (ar *actionResolver) Resolve(cancel <-chan bool) (*Result, error) {
	resp, body, err := ar.exhaust(cancel)
	if err != nil {
		log.Infof("Failed getting response from %v: %v", ar.endpoint, err)
		return nil, nil
	}

	hdlr, ok := handlers[resp.StatusCode]
	if !ok {
		log.Infof("Response %s from %v unusable", resp.Status, ar.endpoint)
		return nil, nil
	}

	return hdlr(resp, body)
}

func (ar *actionResolver) exhaust(cancel <-chan bool) (resp *http.Response, body []byte, err error) {
	var req *http.Request

	req, err = ar.first()
	if err != nil {
		return nil, nil, err
	}

	for req != nil {
		resp, body, err = ar.one(req, cancel)
		if err != nil {
			return nil, nil, err
		}

		req, err = ar.next(resp)
		if err != nil {
			return nil, nil, err
		}
	}

	return resp, body, err
}

func (ar *actionResolver) first() (*http.Request, error) {
	req, err := ar.action.HTTPRequest()
	if err != nil {
		// the inability to build an HTTP request is not recoverable
		return nil, err
	}

	// the URL in the http.Request must not be completely overwritten
	req.URL.Scheme = ar.endpoint.Scheme
	req.URL.Host = ar.endpoint.Host

	return req, nil
}

func (ar *actionResolver) next(resp *http.Response) (*http.Request, error) {
	if resp.StatusCode != http.StatusTemporaryRedirect {
		return nil, nil
	}

	ar.redirectCount += 1
	if ar.redirectCount >= redirectMax {
		return nil, errors.New("too many redirects")
	}

	loc, err := resp.Location()
	if err != nil {
		return nil, err
	}

	req, err := ar.action.HTTPRequest()
	if err != nil {
		return nil, err
	}

	req.URL = loc
	return req, nil
}

func (ar *actionResolver) one(req *http.Request, cancel <-chan bool) (resp *http.Response, body []byte, err error) {
	log.V(1).Infof("etcd: sending HTTP request %s %s", req.Method, req.URL)
	resp, body, err = ar.requestFunc(req, cancel)
	if err != nil {
		log.V(1).Infof("etcd: recv error response from %s %s: %v", req.Method, req.URL, err)
		return
	}

	log.V(1).Infof("etcd: recv response from %s %s: %s", req.Method, req.URL, resp.Status)
	return
}
