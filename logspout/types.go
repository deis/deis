package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
)

type AttachEvent struct {
	Type string
	ID   string
	Name string
}

type Log struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Data string `json:"data"`
}

type Route struct {
	ID     string  `json:"id"`
	Source *Source `json:"source,omitempty"`
	Target Target  `json:"target"`
	closer chan bool
}

type Source struct {
	ID     string   `json:"id,omitempty"`
	Name   string   `json:"name,omitempty"`
	Filter string   `json:"filter,omitempty"`
	Types  []string `json:"types,omitempty"`
}

func (s *Source) All() bool {
	return s.ID == "" && s.Name == "" && s.Filter == ""
}

type Target struct {
	Type      string `json:"type"`
	Addr      string `json:"addr"`
	Protocol  string `json:"protocol"`
	AppendTag string `json:"append_tag,omitempty"`
}

func marshal(obj interface{}) []byte {
	bytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Println("marshal:", err)
	}
	return bytes
}

func unmarshal(input io.ReadCloser, obj interface{}) error {
	body, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}
	return nil
}
