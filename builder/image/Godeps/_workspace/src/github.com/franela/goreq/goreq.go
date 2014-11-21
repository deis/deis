package goreq

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type Request struct {
	headers           []headerTuple
	Method            string
	Uri               string
	Body              interface{}
	QueryString       interface{}
	Timeout           time.Duration
	ContentType       string
	Accept            string
	Host              string
	UserAgent         string
	Insecure          bool
	MaxRedirects      int
	Proxy             string
	Compression       *compression
	BasicAuthUsername string
	BasicAuthPassword string
}

type compression struct {
	writer          func(buffer io.Writer) (io.WriteCloser, error)
	reader          func(buffer io.Reader) (io.ReadCloser, error)
	ContentEncoding string
}

type Response struct {
	StatusCode    int
	ContentLength int64
	Body          *Body
	Header        http.Header
}

type headerTuple struct {
	name  string
	value string
}

type Body struct {
	reader           io.ReadCloser
	compressedReader io.ReadCloser
}

type Error struct {
	timeout bool
	Err     error
}

func (e *Error) Timeout() bool {
	return e.timeout
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (b *Body) Read(p []byte) (int, error) {
	if b.compressedReader != nil {
		return b.compressedReader.Read(p)
	}
	return b.reader.Read(p)
}

func (b *Body) Close() error {
	err := b.reader.Close()
	if b.compressedReader != nil {
		return b.compressedReader.Close()
	}
	return err
}

func (b *Body) FromJsonTo(o interface{}) error {
	if body, err := ioutil.ReadAll(b); err != nil {
		return err
	} else if err := json.Unmarshal(body, o); err != nil {
		return err
	}

	return nil
}

func (b *Body) ToString() (string, error) {
	body, err := ioutil.ReadAll(b)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func Gzip() *compression {
	reader := func(buffer io.Reader) (io.ReadCloser, error) {
		return gzip.NewReader(buffer)
	}
	writer := func(buffer io.Writer) (io.WriteCloser, error) {
		return gzip.NewWriter(buffer), nil
	}
	return &compression{writer: writer, reader: reader, ContentEncoding: "gzip"}
}

func Deflate() *compression {
	reader := func(buffer io.Reader) (io.ReadCloser, error) {
		return flate.NewReader(buffer), nil
	}
	writer := func(buffer io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(buffer, -1)
	}
	return &compression{writer: writer, reader: reader, ContentEncoding: "deflate"}
}

func Zlib() *compression {
	reader := func(buffer io.Reader) (io.ReadCloser, error) {
		return zlib.NewReader(buffer)
	}
	writer := func(buffer io.Writer) (io.WriteCloser, error) {
		return zlib.NewWriter(buffer), nil
	}
	return &compression{writer: writer, reader: reader, ContentEncoding: "deflate"}
}

func paramParse(query interface{}) (string, error) {
	var (
		v = &url.Values{}
		s = reflect.ValueOf(query)
		t = reflect.TypeOf(query)
	)

	switch query.(type) {
	case url.Values:
		return query.(url.Values).Encode(), nil
	default:
		for i := 0; i < s.NumField(); i++ {
			v.Add(strings.ToLower(t.Field(i).Name), fmt.Sprintf("%v", s.Field(i).Interface()))
		}
		return v.Encode(), nil
	}
}

func prepareRequestBody(b interface{}) (io.Reader, error) {
	switch b.(type) {
	case string:
		// treat is as text
		return strings.NewReader(b.(string)), nil
	case io.Reader:
		// treat is as text
		return b.(io.Reader), nil
	case []byte:
		//treat as byte array
		return bytes.NewReader(b.([]byte)), nil
	case nil:
		return nil, nil
	default:
		// try to jsonify it
		j, err := json.Marshal(b)
		if err == nil {
			return bytes.NewReader(j), nil
		}
		return nil, err
	}
}

var defaultDialer = &net.Dialer{Timeout: 1000 * time.Millisecond}
var defaultTransport = &http.Transport{Dial: defaultDialer.Dial, Proxy: http.ProxyFromEnvironment}
var defaultClient = &http.Client{Transport: defaultTransport}
var proxyTransport *http.Transport
var proxyClient *http.Client

func SetConnectTimeout(duration time.Duration) {
	defaultDialer.Timeout = duration
}

func (r *Request) AddHeader(name string, value string) {
	if r.headers == nil {
		r.headers = []headerTuple{}
	}
	r.headers = append(r.headers, headerTuple{name: name, value: value})
}

func (r Request) Do() (*Response, error) {
	var req *http.Request
	var er error
	var transport = defaultTransport
	var client = defaultClient

	if r.Proxy != "" {
		proxyUrl, err := url.Parse(r.Proxy)
		if err != nil {
			// proxy address is in a wrong format
			return nil, &Error{Err: err}
		}
		if proxyTransport == nil {
			proxyTransport = &http.Transport{Dial: defaultDialer.Dial, Proxy: http.ProxyURL(proxyUrl)}
			proxyClient = &http.Client{Transport: proxyTransport}
		} else {
			proxyTransport.Proxy = http.ProxyURL(proxyUrl)
		}
		transport = proxyTransport
		client = proxyClient
	}

	if r.Insecure {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else if transport.TLSClientConfig != nil {
		// the default TLS client (when transport.TLSClientConfig==nil) is
		// already set to verify, so do nothing in that case
		transport.TLSClientConfig.InsecureSkipVerify = false
	}

	b, e := prepareRequestBody(r.Body)
	if e != nil {
		// there was a problem marshaling the body
		return nil, &Error{Err: e}
	}

	if r.QueryString != nil {
		param, e := paramParse(r.QueryString)
		if e != nil {
			return nil, &Error{Err: e}
		}
		r.Uri = r.Uri + "?" + param
	}

	var bodyReader io.Reader
	if b != nil && r.Compression != nil {
		buffer := bytes.NewBuffer([]byte{})
		readBuffer := bufio.NewReader(b)
		writer, err := r.Compression.writer(buffer)
		if err != nil {
			return nil, &Error{Err: err}
		}
		_, e = readBuffer.WriteTo(writer)
		writer.Close()
		if e != nil {
			return nil, &Error{Err: e}
		}
		bodyReader = buffer
	} else {
		bodyReader = b
	}
	req, er = http.NewRequest(r.Method, r.Uri, bodyReader)

	if er != nil {
		// we couldn't parse the URL.
		return nil, &Error{Err: er}
	}

	// add headers to the request
	req.Host = r.Host
	req.Header.Add("User-Agent", r.UserAgent)
	req.Header.Add("Content-Type", r.ContentType)
	req.Header.Add("Accept", r.Accept)
	if r.Compression != nil {
		req.Header.Add("Content-Encoding", r.Compression.ContentEncoding)
		req.Header.Add("Accept-Encoding", r.Compression.ContentEncoding)
	}
	if r.headers != nil {
		for _, header := range r.headers {
			req.Header.Add(header.name, header.value)
		}
	}

	//use basic auth if required
	if r.BasicAuthUsername != "" {
		req.SetBasicAuth(r.BasicAuthUsername, r.BasicAuthPassword)
	}

	timeout := false
	var timer *time.Timer
	if r.Timeout > 0 {
		timer = time.AfterFunc(r.Timeout, func() {
			transport.CancelRequest(req)
			timeout = true
		})
	}

	res, err := client.Do(req)
	if timer != nil {
		timer.Stop()
	}

	if err != nil {
		if !timeout {
			switch err := err.(type) {
			case *net.OpError:
				timeout = err.Timeout()
			case *url.Error:
				if op, ok := err.Err.(*net.OpError); ok {
					timeout = op.Timeout()
				}
			}
		}

		return nil, &Error{timeout: timeout, Err: err}
	}

	if isRedirect(res.StatusCode) && r.MaxRedirects > 0 {
		loc, _ := res.Location()
		r.MaxRedirects--
		r.Uri = loc.String()
		return r.Do()
	}

	if r.Compression != nil && strings.Contains(res.Header.Get("Content-Encoding"), r.Compression.ContentEncoding) {
		compressedReader, err := r.Compression.reader(res.Body)
		if err != nil {
			return nil, &Error{Err: err}
		}
		return &Response{StatusCode: res.StatusCode, ContentLength: res.ContentLength, Header: res.Header, Body: &Body{reader: res.Body, compressedReader: compressedReader}}, nil
	} else {
		return &Response{StatusCode: res.StatusCode, ContentLength: res.ContentLength, Header: res.Header, Body: &Body{reader: res.Body}}, nil
	}
}

func isRedirect(status int) bool {
	switch status {
	case http.StatusMovedPermanently:
		return true
	case http.StatusFound:
		return true
	case http.StatusSeeOther:
		return true
	case http.StatusTemporaryRedirect:
		return true
	default:
		return false
	}
}
