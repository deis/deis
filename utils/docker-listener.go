package main

import (
	"github.com/fsouza/go-dockerclient"
	"log"
)

func assert(err error, context string) {
	if err != nil {
		log.Fatal(context+": ", err)
	}
}

func main() {
	client, err := docker.NewClient("unix:///var/run/docker.sock")
	assert(err, "docker")
	events := make(chan *docker.APIEvents)
	//assert(client.AddEventListener(events), "attacher")
	//assert(client.RemoveEventListener(events), "attacher")
	assert(client.AddEventListener(events), "attacher")
	log.Println("listening for events")
	for msg := range events {
		log.Println("event:", msg.ID[:12], msg.Status)
	}
}
