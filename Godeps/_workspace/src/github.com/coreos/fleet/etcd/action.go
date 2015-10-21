// Copyright 2014 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcd

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

type Action interface {
	fmt.Stringer
	HTTPRequest() (*http.Request, error)
}

//TODO(bcwaldon): Should Delete be separate from CompareAndDelete?
type Delete struct {
	Key           string
	Recursive     bool
	PreviousValue string
	PreviousIndex uint64
}

func (del *Delete) String() string {
	return fmt.Sprintf("{Delete %s}", del.Key)
}

func (del *Delete) HTTPRequest() (*http.Request, error) {
	endpoint := v2URL(del.Key)

	params := endpoint.Query()
	params.Add("recursive", strconv.FormatBool(del.Recursive))
	if del.PreviousValue != "" {
		params.Add("prevValue", del.PreviousValue)
	}
	if del.PreviousIndex != 0 {
		params.Add("prevIndex", strconv.FormatInt(int64(del.PreviousIndex), 10))
	}
	endpoint.RawQuery = params.Encode()

	return http.NewRequest("DELETE", endpoint.String(), nil)
}

type Create struct {
	Key   string
	Value string
	TTL   time.Duration
}

func (c *Create) String() string {
	return fmt.Sprintf("{Create %s}", c.Key)
}

func (c *Create) HTTPRequest() (*http.Request, error) {
	endpoint := v2URL(c.Key)

	params := endpoint.Query()
	params.Add("prevExist", "false")
	endpoint.RawQuery = params.Encode()

	form := url.Values{}
	form.Set("value", c.Value)

	ttl := uint64(c.TTL.Seconds())
	if ttl > 0 {
		form.Set("ttl", strconv.FormatInt(int64(ttl), 10))
	}

	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("PUT", endpoint.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	return req, nil
}

type Update struct {
	Key   string
	Value string
	TTL   time.Duration
}

func (u *Update) String() string {
	return fmt.Sprintf("{Update %s}", u.Key)
}

func (u *Update) HTTPRequest() (*http.Request, error) {
	endpoint := v2URL(u.Key)

	params := endpoint.Query()
	params.Add("prevExist", "true")
	endpoint.RawQuery = params.Encode()

	form := url.Values{}
	form.Set("value", u.Value)

	ttl := int64(u.TTL.Seconds())
	if ttl > 0 {
		form.Set("ttl", strconv.FormatInt(ttl, 10))
	}

	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("PUT", endpoint.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	return req, nil
}

type Set struct {
	Key           string
	Value         string
	TTL           time.Duration
	PreviousIndex uint64
	PreviousValue string
}

func (s *Set) String() string {
	return fmt.Sprintf("{Set %s}", s.Key)
}

func (s *Set) HTTPRequest() (*http.Request, error) {
	endpoint := v2URL(s.Key)

	params := endpoint.Query()
	if s.PreviousIndex != 0 {
		params.Add("prevIndex", strconv.FormatInt(int64(s.PreviousIndex), 10))
	}
	if s.PreviousValue != "" {
		params.Add("prevValue", s.PreviousValue)
	}
	endpoint.RawQuery = params.Encode()

	form := url.Values{}
	form.Set("value", s.Value)

	ttl := int64(s.TTL.Seconds())
	if ttl > 0 {
		form.Set("ttl", strconv.FormatInt(ttl, 10))
	}

	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("PUT", endpoint.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	return req, nil
}

type Get struct {
	Key       string
	Sorted    bool
	Recursive bool
}

func (g *Get) String() string {
	return fmt.Sprintf("{Get %s}", g.Key)
}

func (g *Get) HTTPRequest() (*http.Request, error) {
	endpoint := v2URL(g.Key)

	params := endpoint.Query()
	params.Add("consistent", "true")
	params.Add("sorted", strconv.FormatBool(g.Sorted))
	params.Add("recursive", strconv.FormatBool(g.Recursive))
	endpoint.RawQuery = params.Encode()

	return http.NewRequest("GET", endpoint.String(), nil)
}

type Watch struct {
	Key       string
	Recursive bool
	WaitIndex uint64
}

func (w *Watch) String() string {
	return fmt.Sprintf("{Watch %s}", w.Key)
}

func (w *Watch) HTTPRequest() (*http.Request, error) {
	endpoint := v2URL(w.Key)

	params := endpoint.Query()
	params.Add("consistent", "true")
	params.Add("wait", "true")
	params.Add("recursive", strconv.FormatBool(w.Recursive))
	if w.WaitIndex > 0 {
		params.Set("waitIndex", strconv.FormatInt(int64(w.WaitIndex), 10))
	}
	endpoint.RawQuery = params.Encode()

	return http.NewRequest("GET", endpoint.String(), nil)
}

// v2URL builds a url.URL with an appropriate Path for the v2 etcd API. The
// url.URL's Path attribute is constructed from the base v2 keys path and
// the provided relative path.
func v2URL(rel string) url.URL {
	return url.URL{Path: path.Join("/v2/keys/", rel)}
}
