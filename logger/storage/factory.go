package storage

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/deis/deis/logger/storage/file"
	"github.com/deis/deis/logger/storage/ringbuffer"
)

// Exported so it can be set by an external agent-- namely main.go, which does some flag parsing.
var LogRoot string

var memoryAdapterRegex *regexp.Regexp

func init() {
	memoryAdapterRegex = regexp.MustCompile(`^memory(?::(\d+))?$`)
}

// NewAdapter returns a pointer to an appropriate implementation of the Adapter interface, as
// determined by the storeageAdapterType string it is passed.
func NewAdapter(storeageAdapterType string) (Adapter, error) {
	if storeageAdapterType == "" || storeageAdapterType == "file" {
		adapter, err := file.NewStorageAdapter(LogRoot)
		if err != nil {
			return nil, err
		}
		return adapter, nil
	}
	match := memoryAdapterRegex.FindStringSubmatch(storeageAdapterType)
	if match == nil {
		return nil, fmt.Errorf("Unrecognized storage adapter type: '%s'", storeageAdapterType)
	}
	sizeStr := match[1]
	if sizeStr == "" {
		sizeStr = "1000"
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return nil, err
	}
	adapter, err := ringbuffer.NewStorageAdapter(size)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}
