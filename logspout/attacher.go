package main

import (
	"bufio"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/fsouza/go-dockerclient"
)

type AttachManager struct {
	sync.Mutex
	attached map[string]*LogPump
	channels map[chan *AttachEvent]struct{}
	client   *docker.Client
}

func NewAttachManager(client *docker.Client) *AttachManager {
	m := &AttachManager{
		attached: make(map[string]*LogPump),
		channels: make(map[chan *AttachEvent]struct{}),
		client:   client,
	}
	containers, err := client.ListContainers(docker.ListContainersOptions{})
	assert(err, "attacher")
	for _, listing := range containers {
		m.attach(listing.ID[:12])
	}
	go func() {
		events := make(chan *docker.APIEvents)
		assert(client.AddEventListener(events), "attacher")
		for msg := range events {
			debug("event:", msg.ID[:12], msg.Status)
			if msg.Status == "start" {
				go m.attach(msg.ID[:12])
			}
		}
		log.Fatal("ruh roh") // todo: loop?
	}()
	return m
}

func (m *AttachManager) attach(id string) {
	container, err := m.client.InspectContainer(id)
	assert(err, "attacher")
	name := container.Name[1:]
	success := make(chan struct{})
	failure := make(chan error)
	outrd, outwr := io.Pipe()
	errrd, errwr := io.Pipe()
	go func() {
		err := m.client.AttachToContainer(docker.AttachToContainerOptions{
			Container:    id,
			OutputStream: outwr,
			ErrorStream:  errwr,
			Stdin:        false,
			Stdout:       true,
			Stderr:       true,
			Stream:       true,
			Success:      success,
		})
		outwr.Close()
		errwr.Close()
		debug("attach:", id, "finished")
		if err != nil {
			close(success)
			failure <- err
		}
		m.send(&AttachEvent{Type: "detach", ID: id, Name: name})
		m.Lock()
		delete(m.attached, id)
		m.Unlock()
	}()
	_, ok := <-success
	if ok {
		m.Lock()
		m.attached[id] = NewLogPump(outrd, errrd, id, name)
		m.Unlock()
		success <- struct{}{}
		m.send(&AttachEvent{ID: id, Name: name, Type: "attach"})
		debug("attach:", id, name, "success")
		return
	}
	debug("attach:", id, "failure:", <-failure)
}

func (m *AttachManager) send(event *AttachEvent) {
	m.Lock()
	defer m.Unlock()
	for ch, _ := range m.channels {
		// TODO: log err after timeout and continue
		ch <- event
	}
}

func (m *AttachManager) addListener(ch chan *AttachEvent) {
	m.Lock()
	defer m.Unlock()
	m.channels[ch] = struct{}{}
	go func() {
		for id, pump := range m.attached {
			ch <- &AttachEvent{ID: id, Name: pump.Name, Type: "attach"}
		}
	}()
}

func (m *AttachManager) removeListener(ch chan *AttachEvent) {
	m.Lock()
	defer m.Unlock()
	delete(m.channels, ch)
}

func (m *AttachManager) Get(id string) *LogPump {
	m.Lock()
	defer m.Unlock()
	return m.attached[id]
}

func (m *AttachManager) Listen(source *Source, logstream chan *Log, closer <-chan bool) {
	if source == nil {
		source = new(Source)
	}
	depth := len(m.attached)
	if depth == 0 {
		depth = 1
	}
	events := make(chan *AttachEvent, depth)
	m.addListener(events)
	defer m.removeListener(events)
	for {
		select {
		case event := <-events:
			if event.Type == "attach" && (source.All() ||
				(source.ID != "" && strings.HasPrefix(event.ID, source.ID)) ||
				(source.Name != "" && event.Name == source.Name) ||
				(source.Filter != "" && strings.Contains(event.Name, source.Filter))) {
				pump := m.Get(event.ID)
				pump.AddListener(logstream)
				defer func() {
					if pump != nil {
						pump.RemoveListener(logstream)
					}
				}()
			} else if source.ID != "" && event.Type == "detach" &&
				strings.HasPrefix(event.ID, source.ID) {
				return
			}
		case <-closer:
			return
		}
	}
}

type LogPump struct {
	sync.Mutex
	ID       string
	Name     string
	channels map[chan *Log]struct{}
}

func NewLogPump(stdout, stderr io.Reader, id, name string) *LogPump {
	obj := &LogPump{
		ID:       id,
		Name:     name,
		channels: make(map[chan *Log]struct{}),
	}
	pump := func(typ string, source io.Reader) {
		buf := bufio.NewReader(source)
		for {
			data, err := buf.ReadBytes('\n')
			if err != nil {
				if err != io.EOF {
					debug("pump:", id, typ+":", err)
				}
				return
			}
			obj.send(&Log{
				Data: strings.TrimSuffix(string(data), "\n"),
				ID:   id,
				Name: name,
				Type: typ,
			})
		}
	}
	go pump("stdout", stdout)
	go pump("stderr", stderr)
	return obj
}

func (o *LogPump) send(log *Log) {
	o.Lock()
	defer o.Unlock()
	for ch, _ := range o.channels {
		// TODO: log err after timeout and continue
		ch <- log
	}
}

func (o *LogPump) AddListener(ch chan *Log) {
	o.Lock()
	defer o.Unlock()
	o.channels[ch] = struct{}{}
}

func (o *LogPump) RemoveListener(ch chan *Log) {
	o.Lock()
	defer o.Unlock()
	delete(o.channels, ch)
}
