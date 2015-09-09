package handler

import (
	"container/ring"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"

	"github.com/deis/deis/logger/drain"
	"github.com/deis/deis/logger/syslog"
)

var (
	logStorage = make(map[string]*ring.Ring)
	ringBufferSize int

	regexpForPost = regexp.MustCompile(`^/([-a-z0-9]+)/.*`)
	regexpForGet = regexp.MustCompile(`^/([-a-z0-9]+)/([0-9]+)/.*`)

	mu sync.Mutex
)

// Add log message to main map with ring byffers by project name
func addToStorage(name string, message string) {
	mu.Lock()
	defer mu.Unlock()
	currentRing, ok := logStorage[name]
	if !ok {
		r := ring.New(ringBufferSize)
		r.Value = message
		logStorage[name] = r
		return
	}
	r := currentRing.Next()
	r.Value = message
	logStorage[name] = r
}

// Get specific amount of log lines for application name from main map with ring buffers
func getFromStorage(name string, lines int) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	currentRing, ok := logStorage[name]
	if !ok {
		return "", fmt.Errorf("Could not find logs for project '%s'", name)
	}
	var data string
	getLine := func(line interface{}) {
		if line == nil || lines <= 0 {
			return
		}
		lines -= 1
		if lines == 0 {
			data += fmt.Sprint(line)
		} else {
			data += fmt.Sprintln(line)
		}
	}

	currentRing.Next().Do(getLine)
	return data, nil
}

func RingBufferHandler(bufferSize int, webServicePort int) *Handler {
	ringBufferSize = bufferSize
	h := Handler{
		BaseHandler: syslog.NewBaseHandler(5, filter, false),
	}
	go h.ringBufferLoop()
	go startWebService(webServicePort)
	fmt.Printf("Web service started on %d port\n", webServicePort)
	return &h
}

// Main loop for this handler, each log line will be sended to drain (if DrainURI specified) and copied to log storage
func (h *Handler) ringBufferLoop() {
	for {
		m := h.Get()
		if m == nil {
			break
		}
		if h.DrainURI != "" {
			drain.SendToDrain(m.String(), h.DrainURI)
		}
		message := m.String()
		projectName, err := getProjectName(message)
		if err != nil {
			projectName = "unknown"
			message = fmt.Sprintln(err)
		}
		addToStorage(projectName, message)
	}
	h.End()
}

// Only one http handler which process requests
func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		match := regexpForPost.FindStringSubmatch(r.RequestURI)

		if match == nil {
			fmt.Fprintf(w, "Could not get application name from url: %s", r.RequestURI)
			return
		}

		r.ParseForm()
		value, ok := r.Form["message"]
		if !ok {
			fmt.Fprintln(w, "Could not read from post request, no 'message' param in POST")
			return
		}
		addToStorage(match[1], value[0])
		return
	}

	match := regexpForGet.FindStringSubmatch(r.RequestURI)

	if match == nil {
		fmt.Fprintf(w, "Could not get application name from url: %s", r.RequestURI)
		return
	}

	logLines, err := strconv.Atoi(match[2])
	if err != nil {
		fmt.Fprintln(w, "Unable to get log lines parameter from request")
		return
	}

	data, err := getFromStorage(match[1], logLines)
	if err != nil {
		fmt.Fprintln(w, err)
	} else {
		fmt.Fprint(w, data)
	}
}

// Start web service which serve controller request for get and post logs
func startWebService(port int) {
	http.HandleFunc("/", httpHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Println(err)
	}
}
