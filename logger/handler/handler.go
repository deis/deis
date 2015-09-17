package handler

import (
	"github.com/deis/deis/logger/handler/ringbuffer"
	"github.com/deis/deis/logger/handler/standard"
)

type LogHandler interface {
	Write(string, string) error
	Read(string, int) ([]string, error)
}

func NewStandardHandler(logRoot string) (LogHandler, error) {
	return standard.NewHandler(logRoot)
}

func NewRingBufferHandler(bufferSize int) (LogHandler, error) {
	return ringbuffer.NewHandler(bufferSize)
}
