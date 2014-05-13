package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RouteStore interface {
	Get(id string) (*Route, error)
	GetAll() ([]*Route, error)
	Add(route *Route) error
	Remove(id string) bool
}

type RouteManager struct {
	sync.Mutex
	persistor RouteStore
	attacher  *AttachManager
	routes    map[string]*Route
}

func NewRouteManager(attacher *AttachManager) *RouteManager {
	return &RouteManager{attacher: attacher, routes: make(map[string]*Route)}
}

func (rm *RouteManager) Load(persistor RouteStore) error {
	routes, err := persistor.GetAll()
	if err != nil {
		return err
	}
	for _, route := range routes {
		rm.Add(route)
	}
	rm.persistor = persistor
	return nil
}

func (rm *RouteManager) Get(id string) (*Route, error) {
	rm.Lock()
	defer rm.Unlock()
	route, ok := rm.routes[id]
	if !ok {
		return nil, os.ErrNotExist
	}
	return route, nil
}

func (rm *RouteManager) GetAll() ([]*Route, error) {
	rm.Lock()
	defer rm.Unlock()
	routes := make([]*Route, 0)
	for _, route := range rm.routes {
		routes = append(routes, route)
	}
	return routes, nil
}

func (rm *RouteManager) Add(route *Route) error {
	rm.Lock()
	defer rm.Unlock()
	if route.ID == "" {
		h := sha1.New()
		io.WriteString(h, strconv.Itoa(int(time.Now().UnixNano())))
		route.ID = fmt.Sprintf("%x", h.Sum(nil))[:12]
	}
	route.closer = make(chan bool)
	rm.routes[route.ID] = route
	types := []string{}
	if route.Source != nil {
		types = append(types, route.Source.Types...)
	}
	go func() {
		logstream := make(chan *Log)
		defer close(logstream)
		go syslogStreamer(route.Target, types, logstream)
		rm.attacher.Listen(route.Source, logstream, route.closer)
	}()
	if rm.persistor != nil {
		if err := rm.persistor.Add(route); err != nil {
			log.Println("persistor:", err)
		}
	}
	return nil
}

func (rm *RouteManager) Remove(id string) bool {
	rm.Lock()
	defer rm.Unlock()
	route, ok := rm.routes[id]
	if ok && route.closer != nil {
		route.closer <- true
	}
	delete(rm.routes, id)
	if rm.persistor != nil {
		rm.persistor.Remove(id)
	}
	return ok
}

type RouteFileStore string

func (fs RouteFileStore) Filename(id string) string {
	return string(fs) + "/" + id + ".json"
}

func (fs RouteFileStore) Get(id string) (*Route, error) {
	file, err := os.Open(fs.Filename(id))
	if err != nil {
		return nil, err
	}
	route := new(Route)
	if err = unmarshal(file, route); err != nil {
		return nil, err
	}
	return route, nil
}

func (fs RouteFileStore) GetAll() ([]*Route, error) {
	files, err := ioutil.ReadDir(string(fs))
	if err != nil {
		return nil, err
	}
	var routes []*Route
	for _, file := range files {
		fileparts := strings.Split(file.Name(), ".")
		if len(fileparts) > 1 && fileparts[1] == "json" {
			route, err := fs.Get(fileparts[0])
			if err == nil {
				routes = append(routes, route)
			}
		}
	}
	return routes, nil
}

func (fs RouteFileStore) Add(route *Route) error {
	return ioutil.WriteFile(fs.Filename(route.ID), marshal(route), 0644)
}

func (fs RouteFileStore) Remove(id string) bool {
	if _, err := os.Stat(fs.Filename(id)); err == nil {
		if err := os.Remove(fs.Filename(id)); err != nil {
			return true
		}
	}
	return false
}
