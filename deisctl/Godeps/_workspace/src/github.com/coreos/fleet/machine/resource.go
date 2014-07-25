package machine

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"runtime"
	"strconv"

	"github.com/coreos/fleet/resource"
)

const (
	memInfoPath = "/proc/meminfo"
)

func readLocalResources() (resource.ResourceTuple, error) {
	var res resource.ResourceTuple
	res.Cores = 100 * runtime.NumCPU()

	// TODO(uwedeportivo): determine disk space

	mem, err := readMeminfo()
	if err != nil {
		return res, err
	}
	res.Memory = mem

	return res, nil
}

// parseMeminfo extracts the total amount of memory
// and returns it in MB.
func parseMeminfo(memstr []byte) (int, error) {
	ss := bufio.NewScanner(bytes.NewBuffer(memstr))
	ss.Split(bufio.ScanWords)
	seenMemToken := false
	mem := 0
	for ss.Scan() {
		token := ss.Text()
		if seenMemToken {
			m, err := strconv.Atoi(token)
			if err != nil {
				return 0, err
			}
			mem = m >> 10
			break
		} else if token == "MemTotal:" {
			seenMemToken = true
		}
	}
	if err := ss.Err(); err != nil {
		return 0, err
	}
	return mem, nil
}

// readMeminfo reads /proc/meminfo and returns
// the total amount of memory in MB available on the system.
func readMeminfo() (int, error) {
	memstr, err := ioutil.ReadFile(memInfoPath)
	if err != nil {
		return 0, err
	}
	return parseMeminfo(memstr)
}
