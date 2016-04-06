package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-etcd/etcd"
	dtime "github.com/deis/deis/pkg/time"
	"github.com/fsouza/go-dockerclient"
	"github.com/go-martini/martini"
	"golang.org/x/net/websocket"
)

const (
	MAX_UDP_MSG_BYTES = 65507
	MAX_TCP_MSG_BYTES = 1048576
)

var debugMode bool

func debug(v ...interface{}) {
	if debugMode {
		log.Println(v...)
	}
}

func assert(err error, context string) {
	if err != nil {
		log.Fatalf("%s: %v", context, err)
	}
}

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}

type Colorizer map[string]int

// returns up to 14 color escape codes (then repeats) for each unique key
func (c Colorizer) Get(key string) string {
	i, exists := c[key]
	if !exists {
		c[key] = len(c)
		i = c[key]
	}
	bright := "1;"
	if i%14 > 6 {
		bright = ""
	}
	return "\x1b[" + bright + "3" + strconv.Itoa(7-(i%7)) + "m"
}

func syslogStreamer(target Target, types []string, logstream chan *Log) {
	typestr := "," + strings.Join(types, ",") + ","
	for logline := range logstream {
		if typestr != ",," && !strings.Contains(typestr, logline.Type) {
			continue
		}
		tag, pid, data := getLogParts(logline)

		// HACK: Go's syslog package hardcodes the log format, so let's send our own message
		data = fmt.Sprintf("%s %s[%s]: %s",
			time.Now().Format(getopt("DATETIME_FORMAT", dtime.DeisDatetimeFormat)),
			tag,
			pid,
			data)

		if strings.EqualFold(target.Protocol, "tcp") {
			addr, err := net.ResolveTCPAddr("tcp", target.Addr)
			assert(err, "syslog")
			conn, err := net.DialTCP("tcp", nil, addr)
			assert(err, "syslog")
			assert(conn.SetWriteBuffer(MAX_TCP_MSG_BYTES), "syslog")
			_, err = fmt.Fprintln(conn, data)
			assert(err, "syslog")
		} else if strings.EqualFold(target.Protocol, "udp") {
			// Truncate the message if it's too long to fit in a single UDP packet.
			// Get the bytes first.  If the string has non-UTF8 chars, the number of
			// bytes might exceed the number of characters and it would be good to
			// know that up front.
			dataBytes := []byte(data)
			if len(dataBytes) > MAX_UDP_MSG_BYTES {
				// Truncate the bytes and add ellipses.
				dataBytes = append(dataBytes[:MAX_UDP_MSG_BYTES-3], "..."...)
			}
			addr, err := net.ResolveUDPAddr("udp", target.Addr)
			assert(err, "syslog")
			conn, err := net.DialUDP("udp", nil, addr)
			assert(err, "syslog")
			assert(conn.SetWriteBuffer(MAX_UDP_MSG_BYTES), "syslog")
			_, err = conn.Write(dataBytes)
			assert(err, "syslog")
		} else {
			assert(fmt.Errorf("%s is not a supported protocol, use either udp or tcp", target.Protocol), "syslog")
		}
	}
}

// getLogParts returns a custom tag and PID for containers that
// match Deis' specific application name format. Otherwise,
// it returns the original name and 1 as the PID.  Additionally,
// it returns log data.  The function is also smart enough to
// detect when a leading tag in the log data represents an attempt
// by the controller to log an application event.
func getLogParts(logline *Log) (string, string, string) {
	// example regex that should match: go_v2.web.1
	match := getMatch(`(^[a-z0-9-]+)_(v[0-9]+)\.([a-z-_]+\.[0-9]+)$`, logline.Name)
	if match != nil {
		return match[1], match[3], logline.Data
	}
	if logline.Name == "deis-controller" {
		data_match := getMatch(`^[A-Z]+ \[([a-z0-9-]+)\]: (.*)`, logline.Data)
		if data_match != nil {
			return data_match[1], "deis-controller", data_match[2]
		}
	}
	return logline.Name, "1", logline.Data
}

func getMatch(regex string, name string) []string {
	r := regexp.MustCompile(regex)
	match := r.FindStringSubmatch(name)
	return match
}

func websocketStreamer(w http.ResponseWriter, req *http.Request, logstream chan *Log, closer chan bool) {
	websocket.Handler(func(conn *websocket.Conn) {
		for logline := range logstream {
			if req.URL.Query().Get("type") != "" && logline.Type != req.URL.Query().Get("type") {
				continue
			}
			_, err := conn.Write(append(marshal(logline), '\n'))
			if err != nil {
				closer <- true
				return
			}
		}
	}).ServeHTTP(w, req)
}

func httpStreamer(w http.ResponseWriter, req *http.Request, logstream chan *Log, multi bool) {
	var colors Colorizer
	var usecolor, usejson bool
	nameWidth := 16
	if req.URL.Query().Get("colors") != "off" {
		colors = make(Colorizer)
		usecolor = true
	}
	if req.Header.Get("Accept") == "application/json" {
		w.Header().Add("Content-Type", "application/json")
		usejson = true
	} else {
		w.Header().Add("Content-Type", "text/plain")
	}
	for logline := range logstream {
		if req.URL.Query().Get("types") != "" && logline.Type != req.URL.Query().Get("types") {
			continue
		}
		if usejson {
			w.Write(append(marshal(logline), '\n'))
		} else {
			if multi {
				if len(logline.Name) > nameWidth {
					nameWidth = len(logline.Name)
				}
				if usecolor {
					w.Write([]byte(fmt.Sprintf(
						"%s%"+strconv.Itoa(nameWidth)+"s|%s\x1b[0m\n",
						colors.Get(logline.Name), logline.Name, logline.Data,
					)))
				} else {
					w.Write([]byte(fmt.Sprintf(
						"%"+strconv.Itoa(nameWidth)+"s|%s\n", logline.Name, logline.Data,
					)))
				}
			} else {
				w.Write(append([]byte(logline.Data), '\n'))
			}
		}
		w.(http.Flusher).Flush()
	}
}

func getEtcdValueOrDefault(c *etcd.Client, key string, defaultValue string) string {
	resp, err := c.Get(key, false, false)
	if err != nil {
		if strings.Contains(fmt.Sprintf("%v", err), "Key not found") {
			return defaultValue
		}
		assert(err, "url")
	}
	return resp.Node.Value
}

func getEtcdRoute(client *etcd.Client) *Route {
	hostResp, err := client.Get("/deis/logs/host", false, false)
	assert(err, "url")
	portResp, err := client.Get("/deis/logs/port", false, false)
	assert(err, "url")
	protocol := getEtcdValueOrDefault(client, "/deis/logs/protocol", "udp")
	host := fmt.Sprintf("%s:%s", hostResp.Node.Value, portResp.Node.Value)
	log.Printf("routing all to %s://%s", protocol, host)
	return &Route{ID: "etcd", Target: Target{Type: "syslog", Addr: host, Protocol: protocol}}
}

func main() {
	runtime.GOMAXPROCS(1)
	debugMode = getopt("DEBUG", "") != ""
	port := getopt("PORT", "8000")
	endpoint := getopt("DOCKER_HOST", "unix:///var/run/docker.sock")
	routespath := getopt("ROUTESPATH", "/var/lib/logspout")

	client, err := docker.NewClient(endpoint)
	assert(err, "docker")
	attacher := NewAttachManager(client)
	router := NewRouteManager(attacher)

	// HACK: if we are connecting to etcd, get the logger's connection
	// details from there
	if etcdHost := os.Getenv("ETCD_HOST"); etcdHost != "" {
		connectionString := []string{"http://" + etcdHost + ":4001"}
		debug("etcd:", connectionString[0])
		etcd := etcd.NewClient(connectionString)
		etcd.SetDialTimeout(3 * time.Second)
		router.Add(getEtcdRoute(etcd))
		go func() {
			for {
				// NOTE(bacongobbler): sleep for a bit before doing the discovery loop again
				time.Sleep(10 * time.Second)
				newRoute := getEtcdRoute(etcd)
				oldRoute, err := router.Get(newRoute.ID)
				// router.Get only returns an error if the route doesn't exist. If it does,
				// then we can skip this check and just add the new route to the routing table
				if err == nil &&
					newRoute.Target.Protocol == oldRoute.Target.Protocol &&
					newRoute.Target.Addr == oldRoute.Target.Addr {
					// NOTE(bacongobbler): the two targets are the same; perform a no-op
					continue
				}
				// NOTE(bacongobbler): this operation is a no-op if the route doesn't exist
				router.Remove(oldRoute.ID)
				router.Add(newRoute)
			}
		}()
	}

	if len(os.Args) > 1 {
		u, err := url.Parse(os.Args[1])
		assert(err, "url")
		log.Println("routing all to " + os.Args[1])
		router.Add(&Route{Target: Target{Type: u.Scheme, Addr: u.Host}})
	}

	if _, err := os.Stat(routespath); err == nil {
		log.Println("loading and persisting routes in " + routespath)
		assert(router.Load(RouteFileStore(routespath)), "persistor")
	}

	m := martini.Classic()

	m.Get("/logs(?:/(?P<predicate>[a-zA-Z]+):(?P<value>.+))?", func(w http.ResponseWriter, req *http.Request, params martini.Params) {
		source := new(Source)
		switch {
		case params["predicate"] == "id" && params["value"] != "":
			source.ID = params["value"][:12]
		case params["predicate"] == "name" && params["value"] != "":
			source.Name = params["value"]
		case params["predicate"] == "filter" && params["value"] != "":
			source.Filter = params["value"]
		}

		if source.ID != "" && attacher.Get(source.ID) == nil {
			http.NotFound(w, req)
			return
		}

		logstream := make(chan *Log)
		defer close(logstream)

		var closer <-chan bool
		if req.Header.Get("Upgrade") == "websocket" {
			closerBi := make(chan bool)
			go websocketStreamer(w, req, logstream, closerBi)
			closer = closerBi
		} else {
			go httpStreamer(w, req, logstream, source.All() || source.Filter != "")
			closer = w.(http.CloseNotifier).CloseNotify()
		}

		attacher.Listen(source, logstream, closer)
	})

	m.Get("/routes", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		routes, _ := router.GetAll()
		w.Write(append(marshal(routes), '\n'))
	})

	m.Post("/routes", func(w http.ResponseWriter, req *http.Request) (int, string) {
		route := new(Route)
		if err := unmarshal(req.Body, route); err != nil {
			return http.StatusBadRequest, "Bad request: " + err.Error()
		}

		// TODO: validate?
		router.Add(route)

		w.Header().Add("Content-Type", "application/json")
		return http.StatusCreated, string(append(marshal(route), '\n'))
	})

	m.Get("/routes/:id", func(w http.ResponseWriter, req *http.Request, params martini.Params) {
		route, _ := router.Get(params["id"])
		if route == nil {
			http.NotFound(w, req)
			return
		}
		w.Write(append(marshal(route), '\n'))
	})

	m.Delete("/routes/:id", func(w http.ResponseWriter, req *http.Request, params martini.Params) {
		if ok := router.Remove(params["id"]); !ok {
			http.NotFound(w, req)
		}
	})

	log.Println("logspout serving http on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, m))
}
